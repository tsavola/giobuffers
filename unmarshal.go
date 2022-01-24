// Copyright (c) 2019 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package giobuffers

//go:generate flatc --go --go-namespace flat gio.fbs

import (
	"bytes"
	"image"
	"image/color"

	"gioui.org/f32"
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

	u.unmarshalOps(gtx, buf)

	m := u.strings[u.slot&1] // Mask avoids bounds check.
	for k := range m {
		delete(m, k)
	}
	u.slot = (u.slot + 1) & 1

	return
}

func (u *Unmarshaler) unmarshalOps(gtx layout.Context, buf *flat.OpNode) {
	table := new(flatbuffers.Table)
	if !buf.Op(table) {
		return
	}

	if b := buf.Previous(new(flat.OpNode)); b != nil {
		u.unmarshalOps(gtx, b)
	}

	switch buf.OpType() {
	case flat.OpMacro:
		var buf flat.OpNode
		buf.Init(table.Bytes, table.Pos)
		u.unmarshalMacroOp(gtx, &buf)

	case flat.OpPaintColor:
		var buf flat.PaintColorOp
		buf.Init(table.Bytes, table.Pos)
		paint.ColorOp{
			Color: unmarshalColorNRGBA(buf.Color(new(flat.ColorNRGBA))),
		}.Add(gtx.Ops)

	case flat.OpPaintImage:
		var buf flat.PaintImageOp
		buf.Init(table.Bytes, table.Pos)
		u.unmarshalPaintImageOp(gtx, buf)

	case flat.OpPaintLinearGradient:
		var buf flat.PaintLinearGradientOp
		buf.Init(table.Bytes, table.Pos)
		paint.LinearGradientOp{
			Stop1:  unmarshalF32Point(buf.Stop1(new(flat.F32Point))),
			Color1: unmarshalColorNRGBA(buf.Color1(new(flat.ColorNRGBA))),
			Stop2:  unmarshalF32Point(buf.Stop2(new(flat.F32Point))),
			Color2: unmarshalColorNRGBA(buf.Color2(new(flat.ColorNRGBA))),
		}.Add(gtx.Ops)

	case flat.OpPaint:
		paint.PaintOp{}.Add(gtx.Ops)

	case flat.OpWidgetLabel:
		var buf flat.WidgetLabelLayout
		buf.Init(table.Bytes, table.Pos)
		u.unmarshalWidgetLabelLayout(gtx, buf)
	}
}

func (u *Unmarshaler) unmarshalMacroOp(gtx layout.Context, buf *flat.OpNode) {
	pos := buf.Table().Pos

	m := op.Record(gtx.Ops)
	u.unmarshalOps(gtx, buf)
	m.Stop()

	u.macros[pos] = m
}

func (u *Unmarshaler) unmarshalPaintImageOp(gtx layout.Context, buf flat.PaintImageOp) {
	table := new(flatbuffers.Table)

	if buf.Src(table) {
		switch buf.SrcType() {
		case flat.ImageImageDecode:
			var buf flat.ImageDecode
			buf.Init(table.Bytes, table.Pos)
			img, _, err := image.Decode(bytes.NewReader(buf.DataBytes()))
			check(err)
			paint.NewImageOp(img).Add(gtx.Ops)

		case flat.ImageImageNRGBA:
			var buf flat.ImageNRGBA
			buf.Init(table.Bytes, table.Pos)
			paint.NewImageOp(&image.NRGBA{
				Pix:    buf.PixBytes(),
				Stride: int(buf.Stride()),
				Rect:   unmarshalImageRectangle(buf.Rect(new(flat.ImageRectangle))),
			}).Add(gtx.Ops)
		}
	}
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

func unmarshalColorNRGBA(buf *flat.ColorNRGBA) (r color.NRGBA) {
	if buf != nil {
		r.R = buf.R()
		r.G = buf.G()
		r.B = buf.B()
		r.A = buf.A()
	}
	return
}

func unmarshalF32Point(buf *flat.F32Point) (r f32.Point) {
	if buf != nil {
		r.X = buf.X()
		r.Y = buf.Y()
	}
	return
}

func unmarshalImagePoint(buf *flat.ImagePoint) (r image.Point) {
	if buf != nil {
		r.X = int(buf.X())
		r.Y = int(buf.Y())
	}
	return
}

func unmarshalImageRectangle(buf *flat.ImageRectangle) (r image.Rectangle) {
	if buf != nil {
		r.Min = unmarshalImagePoint(buf.Min(new(flat.ImagePoint)))
		r.Max = unmarshalImagePoint(buf.Max(new(flat.ImagePoint)))
	}
	return
}

func check(err error) { pan.Check(err) }
