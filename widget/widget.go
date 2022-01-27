// Copyright (c) 2022 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package widget

import (
	"image/color"

	"gioui.org/unit"
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/tsavola/giobuffers/flat"
	internal "github.com/tsavola/giobuffers/internal/resource"
	"github.com/tsavola/giobuffers/layout"
	"github.com/tsavola/giobuffers/op"
	"github.com/tsavola/giobuffers/op/paint"
	"github.com/tsavola/giobuffers/resource"
	"github.com/tsavola/giobuffers/text"
)

type Fit = flat.WidgetFit

const (
	Unscaled  = flat.WidgetFitUnscaled
	Contain   = flat.WidgetFitContain
	Cover     = flat.WidgetFitCover
	ScaleDown = flat.WidgetFitScaleDown
	Fill      = flat.WidgetFitFill
)

type Border struct {
	Color        color.NRGBA
	CornerRadius unit.Value
	Width        unit.Value
}

func (x Border) create(b *flatbuffers.Builder) flatbuffers.UOffsetT {
	flat.WidgetBorderStart(b)
	flat.WidgetBorderAddColor(b, flat.CreateColorNRGBA(b, x.Color.R, x.Color.G, x.Color.B, x.Color.A))
	if x.CornerRadius != (unit.Value{}) {
		flat.WidgetBorderAddCornerRadius(b, flat.CreateUnitValue(b, x.CornerRadius.V, flat.Unit(x.CornerRadius.U)))
	}
	flat.WidgetBorderAddWidth(b, flat.CreateUnitValue(b, x.Width.V, flat.Unit(x.Width.U)))
	return flat.WidgetBorderEnd(b)
}

func (x Border) CreateLayout(b *flatbuffers.Builder, widget op.Ops) op.Op {
	borderOff := x.create(b)

	flat.WidgetBorderLayoutStart(b)
	flat.WidgetBorderLayoutAddBorder(b, borderOff)
	flat.WidgetBorderLayoutAddWidget(b, widget.LastNode)
	offset := flat.WidgetBorderLayoutEnd(b)

	return op.Op{
		Type:   flat.OpWidgetBorder,
		Offset: offset,
	}
}

func (x Border) AddLayout(b *flatbuffers.Builder, o *op.Ops, widget op.Ops) op.Ops {
	return x.CreateLayout(b, widget).Add(b, o)
}

type Icon struct {
	resource.Resource
	data []byte
}

func NewIcon(c *resource.Context, id string, data []byte) *Icon {
	return &Icon{
		internal.MakeResource(c, id),
		data,
	}
}

func (x Icon) create(b *flatbuffers.Builder) flatbuffers.UOffsetT {
	dataOff := b.CreateByteVector(x.data)

	flat.WidgetIconStart(b)
	flat.WidgetIconAddData(b, dataOff)
	return flat.WidgetIconEnd(b)
}

func (x Icon) CreateLayout(b *flatbuffers.Builder, color color.NRGBA) op.Op {
	idOff := b.CreateSharedString(x.ID)

	var iconOff flatbuffers.UOffsetT
	if x.NeedSync() {
		iconOff = x.create(b)
	}

	flat.WidgetIconLayoutStart(b)
	flat.WidgetIconLayoutAddIconId(b, idOff)
	if iconOff != 0 {
		flat.WidgetIconLayoutAddIcon(b, iconOff)
	}
	flat.WidgetIconLayoutAddColor(b, flat.CreateColorNRGBA(b, color.R, color.G, color.B, color.A))
	offset := flat.WidgetIconLayoutEnd(b)

	x.Synced()

	return op.Op{
		Type:   flat.OpWidgetIcon,
		Offset: offset,
	}
}

func (x Icon) AddLayout(b *flatbuffers.Builder, o *op.Ops, color color.NRGBA) op.Ops {
	return x.CreateLayout(b, color).Add(b, o)
}

type Image struct {
	Src      paint.ImageOp
	Fit      Fit
	Position layout.Direction
	Scale    float32
}

func (x Image) create(b *flatbuffers.Builder) flatbuffers.UOffsetT {
	op := x.Src.Create(b)

	flat.WidgetImageStart(b)
	flat.WidgetImageAddSrc(b, op.Offset)
	flat.WidgetImageAddFit(b, x.Fit)
	flat.WidgetImageAddPosition(b, flat.LayoutDirection(x.Position))
	flat.WidgetImageAddScale(b, x.Scale)
	return flat.WidgetImageEnd(b)
}

func (x Image) CreateLayout(b *flatbuffers.Builder) op.Op {
	imageOff := x.create(b)

	flat.WidgetImageLayoutStart(b)
	flat.WidgetImageLayoutAddImage(b, imageOff)
	offset := flat.WidgetImageLayoutEnd(b)

	return op.Op{
		Type:   flat.OpWidgetImage,
		Offset: offset,
	}
}

func (x Image) AddLayout(b *flatbuffers.Builder, o *op.Ops) op.Ops {
	return x.CreateLayout(b).Add(b, o)
}

type Label struct {
	Alignment text.Alignment
	MaxLines  int32
}

func (x Label) CreateLayout(b *flatbuffers.Builder, font text.Font, size unit.Value, txt string) op.Op {
	flat.WidgetLabelStart(b)
	if x.Alignment != 0 {
		flat.WidgetLabelAddAlignment(b, x.Alignment)
	}
	if x.MaxLines != 0 {
		flat.WidgetLabelAddMaxLines(b, int32(x.MaxLines))
	}
	labelOff := flat.WidgetLabelEnd(b)

	variantOff := b.CreateSharedString(font.Variant)

	flat.TextFontStart(b)
	flat.TextFontAddVariant(b, variantOff)
	if font.Weight != 0 {
		flat.TextFontAddWeight(b, int32(font.Weight))
	}
	fontOff := flat.TextFontEnd(b)

	txtOff := b.CreateSharedString(txt)

	flat.WidgetLabelLayoutStart(b)
	flat.WidgetLabelLayoutAddLabel(b, labelOff)
	flat.WidgetLabelLayoutAddFont(b, fontOff)
	flat.WidgetLabelLayoutAddSize(b, flat.CreateUnitValue(b, size.V, flat.Unit(size.U)))
	flat.WidgetLabelLayoutAddText(b, txtOff)
	offset := flat.WidgetLabelLayoutEnd(b)

	return op.Op{
		Type:   flat.OpWidgetLabel,
		Offset: offset,
	}
}

func (x Label) AddLayout(b *flatbuffers.Builder, o *op.Ops, font text.Font, size unit.Value, txt string) op.Ops {
	return x.CreateLayout(b, font, size, txt).Add(b, o)
}
