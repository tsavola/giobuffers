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

type macroMap map[flatbuffers.UOffsetT]ui.MacroOp

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

func Unmarshal(data []byte, ops *ui.Ops, faces *measure.Faces) (err error) {
	defer func() {
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

	unmarshalOps(*buf, ops, faces, make(macroMap))
	return
}

func unmarshalOps(buf flat.OpNode, ops *ui.Ops, faces *measure.Faces, macros macroMap) {
	table := new(flatbuffers.Table)

	for {
		if buf.Op(table) {
			switch buf.OpType() {
			case flat.OpMacro:
				var x flat.OpNode
				x.Init(table.Bytes, table.Pos)
				unmarshalMacroOp(x, ops, faces, macros)

			case flat.OpColor:
				var x flat.ColorOp
				x.Init(table.Bytes, table.Pos)
				unmarshalColorOp(x, ops)

			case flat.OpLabelLayout:
				var x flat.LabelLayout
				x.Init(table.Bytes, table.Pos)
				unmarshalLabelLayout(x, ops, faces, macros)
			}
		}

		if buf.Next(&buf) == nil {
			break
		}
	}
}

func unmarshalMacroOp(buf flat.OpNode, ops *ui.Ops, faces *measure.Faces, macros macroMap) {
	pos := buf.Table().Pos

	var m ui.MacroOp
	m.Record(ops)
	unmarshalOps(buf, ops, faces, macros)
	m.Stop()

	macros[pos] = m
}

func unmarshalColorOp(buf flat.ColorOp, ops *ui.Ops) {
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

func unmarshalLabelLayout(buf flat.LabelLayout, ops *ui.Ops, faces *measure.Faces, macros macroMap) {
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
			m, found := macros[pos]
			if !found {
				panic(fmt.Errorf("material macro not found (%d)", pos))
			}
			l.Material = m
		}

		l.Alignment = text.Alignment(bufL.Alignment())
		l.Text = string(bufL.Text())
		l.MaxLines = int(bufL.MaxLines())

		l.Layout(ops, cs)
	}
}
