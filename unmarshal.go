// Copyright (c) 2019 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package giobuffers

//go:generate flatc --go --go-namespace flat gio.fbs

import (
	//"fmt"
	"image/color"

	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/tsavola/giobuffers/flat"
	"import.name/pan"
)

type Unmarshaler struct {
	shaper  text.Shaper
	macros  map[flatbuffers.UOffsetT]op.MacroOp
	strings [2]map[string]string // Caches strings between two Unmarshal calls.
	slot    int                  // Current string map for lookup.
}

func (u *Unmarshaler) Unmarshal(gtx layout.Context, data []byte) (err error) {
	if u.shaper == nil {
		u.shaper = text.NewCache(gofont.Collection())
		u.macros = make(map[flatbuffers.UOffsetT]op.MacroOp)
		u.strings = [2]map[string]string{
			make(map[string]string),
			make(map[string]string),
		}
	}

	defer func() {
		for k := range u.macros {
			delete(u.macros, k)
		}
	}()

	defer func() {
		err = pan.Error(recover())
	}()

	buf := flat.GetRootAsOpNode(data, 0)
	if buf == nil {
		return
	}

	u.unmarshalOps(gtx, *buf)

	m := u.strings[u.slot&1] // Mask avoids bounds check.
	for k := range m {
		delete(m, k)
	}
	u.slot = (u.slot + 1) & 1

	return
}

func (u *Unmarshaler) unmarshalOps(gtx layout.Context, buf flat.OpNode) {
	table := new(flatbuffers.Table)

	for {
		if buf.Op(table) {
			switch buf.OpType() {
			case flat.OpMacro:
				var x flat.OpNode
				x.Init(table.Bytes, table.Pos)
				u.unmarshalMacroOp(gtx, x)

			case flat.OpPaintColor:
				var x flat.PaintColorOp
				x.Init(table.Bytes, table.Pos)
				u.unmarshalPaintColorOp(gtx, x)

			case flat.OpPaint:
				paint.PaintOp{}.Add(gtx.Ops)

			case flat.OpWidgetLabel:
				var x flat.WidgetLabelLayout
				x.Init(table.Bytes, table.Pos)
				u.unmarshalWidgetLabelLayout(gtx, x)
			}
		}

		if buf.Next(&buf) == nil {
			break
		}
	}
}

func (u *Unmarshaler) unmarshalMacroOp(gtx layout.Context, buf flat.OpNode) {
	pos := buf.Table().Pos

	m := op.Record(gtx.Ops)
	u.unmarshalOps(gtx, buf)
	m.Stop()

	u.macros[pos] = m
}

func (u *Unmarshaler) unmarshalPaintColorOp(gtx layout.Context, buf flat.PaintColorOp) {
	var co paint.ColorOp

	if c := buf.Color(new(flat.ColorNRGBA)); c != nil {
		co.Color = color.NRGBA{
			R: c.R(),
			G: c.G(),
			B: c.B(),
			A: c.A(),
		}
	}

	co.Add(gtx.Ops)
}

func (u *Unmarshaler) unmarshalWidgetLabelLayout(gtx layout.Context, buf flat.WidgetLabelLayout) {
	var label widget.Label
	if b := buf.Label(new(flat.WidgetLabel)); b != nil {
		label.Alignment = text.Alignment(b.Alignment())
		label.MaxLines = int(b.MaxLines())
	}

	var font text.Font
	if b := buf.Font(new(flat.TextFont)); b != nil {
		font.Typeface = text.Typeface(u.unmarshalString(b.Typeface()))
		font.Variant = text.Variant(u.unmarshalString(b.Variant()))
		font.Style = text.Style(b.Style())
		font.Weight = text.Weight(b.Weight())
	}

	var size unit.Value
	if b := buf.Size(new(flat.UnitValue)); b != nil {
		size.V = b.V()
		size.U = unit.Unit(b.U())
	}

	label.Layout(gtx, u.shaper, font, size, u.unmarshalString(buf.Text()))
}

func (u *Unmarshaler) unmarshalString(buf []byte) (s string) {
	// Repeat string(buf) so that the map lookup gets optimized.  The redundant
	// index mask avoids bounds check.
	if existing, found := u.strings[u.slot&1][string(buf)]; found {
		s = existing
	} else {
		s = string(buf)
	}
	u.strings[(u.slot+1)&1][s] = s // Cache it for next round.
	return
}
