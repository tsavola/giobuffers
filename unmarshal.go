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

type cachedImage struct {
	id    string
	image image.Image
}

type cachedIcon struct {
	id   string
	icon *widget.Icon
}

type cache struct {
	strings map[string]string
	images  map[string]cachedImage
	icons   map[string]cachedIcon
}

func makeCache() cache {
	return cache{
		strings: make(map[string]string),
		images:  make(map[string]cachedImage),
		icons:   make(map[string]cachedIcon),
	}
}

func (c *cache) clear() {
	clear(c.strings)
	clear(c.images)
	clear(c.icons)
}

type Unmarshaler struct {
	shaper text.Shaper
	macros map[flatbuffers.UOffsetT]op.MacroOp
	older  cache // From previous iteration.
	newer  cache // For current/next iteration.
}

func (u *Unmarshaler) Unmarshal(gtx layout.Context, buf []byte) (d layout.Dimensions, err error) {
	if u.shaper == nil {
		u.shaper = text.NewCache(gofont.Collection())
		u.macros = make(map[flatbuffers.UOffsetT]op.MacroOp)
		u.older = makeCache()
		u.newer = makeCache()
	}

	defer func() {
		clear(u.macros)

		u.older, u.newer = u.newer, u.older
		u.newer.clear()
	}()

	defer func() {
		err = pan.Error(recover())
	}()

	if b := flat.GetRootAsOpNode(buf, 0); b != nil {
		d = u.ops(gtx, *b)
	}
	return
}

func (u *Unmarshaler) ops(gtx layout.Context, buf flat.OpNode) (d layout.Dimensions) {
	if table := new(flatbuffers.Table); buf.Op(table) {
		if b := buf.Previous(nil); b != nil {
			d = u.ops(gtx, *b)
		}

		switch buf.OpType() {
		case flat.OpMacro:
			var b flat.OpNode
			b.Init(table.Bytes, table.Pos)

			u.macroOp(gtx, b)

		case flat.OpLayoutDirection:
			var b flat.LayoutDirectionLayout
			b.Init(table.Bytes, table.Pos)

			d = layout.Direction(b.Direction()).Layout(gtx, u.layoutWidget(b.Widget(nil)))

		case flat.OpLayoutFlex:
			var b flat.LayoutFlexLayout
			b.Init(table.Bytes, table.Pos)

			d = u.layoutFlexLayout(gtx, b)

		case flat.OpLayoutInset:
			var b flat.LayoutInsetLayout
			b.Init(table.Bytes, table.Pos)

			d = u.layoutInsetLayout(gtx, b)

		case flat.OpLayoutSpacer:
			var b flat.LayoutSpacerLayout
			b.Init(table.Bytes, table.Pos)

			d = u.layoutSpacerLayout(gtx, b)

		case flat.OpLayoutStack:
			var b flat.LayoutStackLayout
			b.Init(table.Bytes, table.Pos)

			d = u.layoutStackLayout(gtx, b)

		case flat.OpPaintColor:
			var b flat.PaintColorOp
			b.Init(table.Bytes, table.Pos)

			paint.ColorOp{
				Color: colorNRGBA(b.Color(nil)),
			}.Add(gtx.Ops)

		case flat.OpPaintImage:
			var b flat.PaintImageOp
			b.Init(table.Bytes, table.Pos)

			if op, ok := u.paintImageOp(gtx, b); ok {
				op.Add(gtx.Ops)
			} else {
				paint.ColorOp{}.Add(gtx.Ops) // Transparent.
			}

		case flat.OpPaintLinearGradient:
			var b flat.PaintLinearGradientOp
			b.Init(table.Bytes, table.Pos)

			paint.LinearGradientOp{
				Stop1:  f32Point(b.Stop1(nil)),
				Color1: colorNRGBA(b.Color1(nil)),
				Stop2:  f32Point(b.Stop2(nil)),
				Color2: colorNRGBA(b.Color2(nil)),
			}.Add(gtx.Ops)

		case flat.OpPaint:
			var b flat.PaintOp
			b.Init(table.Bytes, table.Pos)

			paint.PaintOp{}.Add(gtx.Ops)

		case flat.OpWidgetBorder:
			var b flat.WidgetBorderLayout
			b.Init(table.Bytes, table.Pos)

			d = u.widgetBorderLayout(gtx, b)

		case flat.OpWidgetIcon:
			var b flat.WidgetIconLayout
			b.Init(table.Bytes, table.Pos)

			d = u.widgetIconLayout(gtx, b)

		case flat.OpWidgetImage:
			var b flat.WidgetImageLayout
			b.Init(table.Bytes, table.Pos)

			d = u.widgetImageLayout(gtx, b)

		case flat.OpWidgetLabel:
			var b flat.WidgetLabelLayout
			b.Init(table.Bytes, table.Pos)

			d = u.widgetLabelLayout(gtx, b)
		}
	}

	return
}

func (u *Unmarshaler) macroOp(gtx layout.Context, buf flat.OpNode) {
	pos := buf.Table().Pos

	m := op.Record(gtx.Ops)
	u.ops(gtx, buf)
	m.Stop()

	u.macros[pos] = m
}

func (u *Unmarshaler) layoutFlexLayout(gtx layout.Context, buf flat.LayoutFlexLayout) layout.Dimensions {
	var flex layout.Flex
	if b := buf.Flex(nil); b != nil {
		flex.Axis = layout.Axis(b.Axis())
		flex.Spacing = layout.Spacing(b.Spacing())
		flex.Alignment = layout.Alignment(b.Alignment())
		flex.WeightSum = b.WeightSum()
	}

	var children []layout.FlexChild
	for b := buf.Children(nil); b != nil; b = b.Next(b) {
		children = append(children, u.layoutFlexChild(gtx, *b))
	}

	return flex.Layout(gtx, children...)
}

func (u *Unmarshaler) layoutFlexChild(gtx layout.Context, buf flat.LayoutFlexChildNode) layout.FlexChild {
	if table := new(flatbuffers.Table); buf.Child(table) {
		switch buf.ChildType() {
		case flat.LayoutFlexChildFlexed:
			var buf flat.LayoutFlexed
			buf.Init(table.Bytes, table.Pos)

			return layout.Flexed(buf.Weight(), u.layoutWidget(buf.Widget(nil)))

		case flat.LayoutFlexChildRigid:
			var buf flat.LayoutRigid
			buf.Init(table.Bytes, table.Pos)

			return layout.Rigid(u.layoutWidget(buf.Widget(nil)))
		}
	}

	return layout.Rigid(u.layoutWidget(nil))
}

func (u *Unmarshaler) layoutInsetLayout(gtx layout.Context, buf flat.LayoutInsetLayout) layout.Dimensions {
	var inset layout.Inset
	if b := buf.Inset(nil); b != nil {
		inset.Top = unitValue(b.Top(nil))
		inset.Bottom = unitValue(b.Bottom(nil))
		inset.Left = unitValue(b.Left(nil))
		inset.Right = unitValue(b.Right(nil))
	}

	return inset.Layout(gtx, u.layoutWidget(buf.Widget(nil)))
}

func (u *Unmarshaler) layoutSpacerLayout(gtx layout.Context, buf flat.LayoutSpacerLayout) layout.Dimensions {
	var spacer layout.Spacer
	if b := buf.Spacer(nil); b != nil {
		spacer.Width = unitValue(b.Width(nil))
		spacer.Height = unitValue(b.Height(nil))
	}

	return spacer.Layout(gtx)
}

func (u *Unmarshaler) layoutStackLayout(gtx layout.Context, buf flat.LayoutStackLayout) layout.Dimensions {
	var stack layout.Stack
	if b := buf.Stack(nil); b != nil {
		stack.Alignment = layout.Direction(b.Alignment())
	}

	var children []layout.StackChild
	for b := buf.Children(nil); b != nil; b = b.Next(b) {
		children = append(children, u.layoutStackChild(gtx, *b))
	}

	return stack.Layout(gtx, children...)
}

func (u *Unmarshaler) layoutStackChild(gtx layout.Context, buf flat.LayoutStackChildNode) layout.StackChild {
	if table := new(flatbuffers.Table); buf.Child(table) {
		switch buf.ChildType() {
		case flat.LayoutStackChildExpanded:
			var buf flat.LayoutExpanded
			buf.Init(table.Bytes, table.Pos)

			return layout.Expanded(u.layoutWidget(buf.Widget(nil)))

		case flat.LayoutStackChildStacked:
			var buf flat.LayoutStacked
			buf.Init(table.Bytes, table.Pos)

			return layout.Stacked(u.layoutWidget(buf.Widget(nil)))
		}
	}

	return layout.Stacked(u.layoutWidget(nil))
}

func (u *Unmarshaler) layoutWidget(buf *flat.OpNode) layout.Widget {
	return func(gtx layout.Context) (d layout.Dimensions) {
		if buf != nil {
			d = u.ops(gtx, *buf)
		}
		return
	}
}

func (u *Unmarshaler) paintImageOp(gtx layout.Context, buf flat.PaintImageOp) (paint.ImageOp, bool) {
	c, existed := u.newer.images[string(buf.SrcId())]
	if !existed {
		c, existed = u.older.images[string(buf.SrcId())]
	}

	if table := new(flatbuffers.Table); buf.Src(table) {
		switch buf.SrcType() {
		case flat.ImageImageDecode:
			var buf flat.ImageDecode
			buf.Init(table.Bytes, table.Pos)
			img, _, err := image.Decode(bytes.NewReader(buf.DataBytes()))
			check(err)
			c.image = img

		case flat.ImageImageNRGBA:
			var buf flat.ImageNRGBA
			buf.Init(table.Bytes, table.Pos)
			c.image = &image.NRGBA{
				Pix:    append([]byte{}, buf.PixBytes()...),
				Stride: int(buf.Stride()),
				Rect:   imageRectangle(buf.Rect(nil)),
			}
		}
	}

	if c.image == nil {
		return paint.ImageOp{}, false
	}

	if !existed {
		c.id = string(buf.SrcId())
	}
	u.newer.images[c.id] = c

	return paint.NewImageOp(c.image), true
}

func (u *Unmarshaler) widgetBorderLayout(gtx layout.Context, buf flat.WidgetBorderLayout) layout.Dimensions {
	var border widget.Border
	if b := buf.Border(nil); b != nil {
		border.Color = colorNRGBA(b.Color(nil))
		border.CornerRadius = unitValue(b.CornerRadius(nil))
		border.Width = unitValue(b.Width(nil))
	}

	return border.Layout(gtx, u.layoutWidget(buf.Widget(nil)))
}

func (u *Unmarshaler) widgetIconLayout(gtx layout.Context, buf flat.WidgetIconLayout) (_ layout.Dimensions) {
	c, existed := u.newer.icons[string(buf.IconId())]
	if !existed {
		c, existed = u.older.icons[string(buf.IconId())]
	}

	if b := buf.Icon(nil); b != nil {
		icon, err := widget.NewIcon(b.DataBytes())
		check(err)
		c.icon = icon
	}

	if c.icon == nil {
		return
	}

	if !existed {
		c.id = string(buf.IconId())
	}
	u.newer.icons[c.id] = c

	return c.icon.Layout(gtx, colorNRGBA(buf.Color(nil)))
}

func (u *Unmarshaler) widgetImageLayout(gtx layout.Context, buf flat.WidgetImageLayout) (_ layout.Dimensions) {
	imageBuf := buf.Image(nil)
	if imageBuf == nil {
		return
	}

	srcBuf := imageBuf.Src(nil)
	if srcBuf == nil {
		return
	}

	src, ok := u.paintImageOp(gtx, *srcBuf)
	if !ok {
		return
	}

	return widget.Image{
		Src:      src,
		Fit:      widget.Fit(imageBuf.Fit()),
		Position: layout.Direction(imageBuf.Position()),
		Scale:    imageBuf.Scale(),
	}.Layout(gtx)
}

func (u *Unmarshaler) widgetLabelLayout(gtx layout.Context, buf flat.WidgetLabelLayout) layout.Dimensions {
	var label widget.Label
	if b := buf.Label(nil); b != nil {
		label.Alignment = text.Alignment(b.Alignment())
		label.MaxLines = int(b.MaxLines())
	}

	var font text.Font
	if b := buf.Font(nil); b != nil {
		font.Typeface = text.Typeface(u.string(b.Typeface()))
		font.Variant = text.Variant(u.string(b.Variant()))
		font.Style = text.Style(b.Style())
		font.Weight = text.Weight(b.Weight())
	}

	var size unit.Value
	if b := buf.Size(nil); b != nil {
		size.V = b.V()
		size.U = unit.Unit(b.U())
	}

	return label.Layout(gtx, u.shaper, font, size, u.string(buf.Text()))
}

func (u *Unmarshaler) string(buf []byte) string {
	// Repeat string(buf) so that the map lookup gets optimized.

	if s, found := u.newer.strings[string(buf)]; found {
		return s
	}

	s, found := u.older.strings[string(buf)]
	if !found {
		s = string(buf)
	}
	u.newer.strings[s] = s // Cache it for next round.
	return s
}

func colorNRGBA(b *flat.ColorNRGBA) (c color.NRGBA) {
	if b != nil {
		c.R = b.R()
		c.G = b.G()
		c.B = b.B()
		c.A = b.A()
	}
	return
}

func f32Point(b *flat.F32Point) (p f32.Point) {
	if b != nil {
		p.X = b.X()
		p.Y = b.Y()
	}
	return
}

func imagePoint(b *flat.ImagePoint) (p image.Point) {
	if b != nil {
		p.X = int(b.X())
		p.Y = int(b.Y())
	}
	return
}

func imageRectangle(b *flat.ImageRectangle) (r image.Rectangle) {
	if b != nil {
		r.Min = imagePoint(b.Min(nil))
		r.Max = imagePoint(b.Max(nil))
	}
	return
}

func unitValue(b *flat.UnitValue) (r unit.Value) {
	if b != nil {
		r.V = b.V()
		r.U = unit.Unit(b.U())
	}
	return
}

func check(err error) { pan.Check(err) }
