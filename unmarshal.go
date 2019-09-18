// Copyright (c) 2019 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package giobuffers

//go:generate flatc --go gio.fbs

import (
	"fmt"
	"image/color"

	"gioui.org/ui"
	"gioui.org/ui/layout"
	"gioui.org/ui/measure"
	"gioui.org/ui/paint"
	"gioui.org/ui/text"
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/tsavola/giobuffers/flat"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/sfnt"
)

var fonts map[string]*sfnt.Font

func init() {
	regular, err := sfnt.Parse(goregular.TTF)
	if err != nil {
		panic("failed to load font")
	}

	fonts = map[string]*sfnt.Font{
		"Go Regular": regular,
	}
}

type Unmarshaler struct {
	macros  map[flatbuffers.UOffsetT]ui.MacroOp
	strings [2]map[string]string // Caches strings between two Unmarshal calls.
	slot    int                  // Current string map for lookup.
}

func (u *Unmarshaler) Unmarshal(data []byte, ops *ui.Ops, faces *measure.Faces) (err error) {
	if u.macros == nil {
		u.macros = make(map[flatbuffers.UOffsetT]ui.MacroOp)
		u.strings = [2]map[string]string{
			make(map[string]string),
			make(map[string]string),
		}
	}

	defer func() {
		for k := range u.macros {
			delete(u.macros, k)
		}

		if x := recover(); x != nil {
			if e, ok := x.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("panic: %v", x)
			}
		}
	}()

	buf := flat.GetRootAsOpNode(data, 0)
	if buf == nil {
		return
	}

	u.unmarshalOps(*buf, ops, faces)

	m := u.strings[u.slot&1] // Mask avoids bounds check.
	for k := range m {
		delete(m, k)
	}
	u.slot = (u.slot + 1) & 1

	return
}

func (u *Unmarshaler) unmarshalOps(buf flat.OpNode, ops *ui.Ops, faces *measure.Faces) {
	table := new(flatbuffers.Table)

	for {
		if buf.Op(table) {
			switch buf.OpType() {
			case flat.OpMacro:
				var x flat.OpNode
				x.Init(table.Bytes, table.Pos)
				u.unmarshalMacroOp(x, ops, faces)

			case flat.OpColor:
				var x flat.ColorOp
				x.Init(table.Bytes, table.Pos)
				u.unmarshalColorOp(x, ops)

			case flat.OpLabelLayout:
				var x flat.LabelLayout
				x.Init(table.Bytes, table.Pos)
				u.unmarshalLabelLayout(x, ops, faces)
			}
		}

		if buf.Next(&buf) == nil {
			break
		}
	}
}

func (u *Unmarshaler) unmarshalMacroOp(buf flat.OpNode, ops *ui.Ops, faces *measure.Faces) {
	pos := buf.Table().Pos

	var m ui.MacroOp
	m.Record(ops)
	u.unmarshalOps(buf, ops, faces)
	m.Stop()

	u.macros[pos] = m
}

func (u *Unmarshaler) unmarshalColorOp(buf flat.ColorOp, ops *ui.Ops) {
	var co paint.ColorOp

	if c := buf.Color(new(flat.ColorRGBA)); c != nil {
		co.Color = color.RGBA{
			R: c.R(),
			G: c.G(),
			B: c.B(),
			A: c.A(),
		}
	}

	co.Add(ops)
}

func (u *Unmarshaler) unmarshalLabelLayout(buf flat.LabelLayout, ops *ui.Ops, faces *measure.Faces) {
	var cs layout.Constraints

	if bufCS := buf.Constraints(new(flat.Constraints)); bufCS != nil {
		c := bufCS.Width(new(flat.Constraint))
		cs.Width = layout.Constraint{
			Min: int(c.Min()),
			Max: int(c.Max()),
		}

		c = bufCS.Height(c)
		cs.Height = layout.Constraint{
			Min: int(c.Min()),
			Max: int(c.Max()),
		}
	}

	if bufL := buf.Label(new(flat.Label)); bufL != nil {
		var l text.Label

		if f := bufL.Face(new(flat.Face)); f != nil {
			name := f.Font()
			font := fonts[string(name)]
			if font == nil {
				panic(fmt.Errorf("unsupported font: %q", name))
			}

			var s ui.Value

			if bufS := f.Size(new(flat.Value)); bufS != nil {
				s = ui.Value{
					V: bufS.V(),
					U: ui.Unit(bufS.U()),
				}
			}

			l.Face = faces.For(font, s)
		}

		if bufM := bufL.Material(new(flat.OpNode)); bufM != nil {
			pos := bufM.Table().Pos
			m, found := u.macros[pos]
			if !found {
				panic(fmt.Errorf("material macro not found (%d)", pos))
			}
			l.Material = m
		}

		l.Alignment = text.Alignment(bufL.Alignment())

		// Repeat string(text) so that the map lookup gets optimized.
		// The redundant index mask avoids bounds check.
		text := bufL.Text()
		s, found := u.strings[u.slot&1][string(text)]
		if !found {
			s = string(text)
		}
		u.strings[(u.slot+1)&1][s] = s // Cache it for next round.
		l.Text = s

		l.MaxLines = int(bufL.MaxLines())

		l.Layout(ops, cs)
	}
}
