// Copyright (c) 2022 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package paint

import (
	"image"
	"image/color"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/tsavola/giobuffers/flat"
	internal "github.com/tsavola/giobuffers/internal/resource"
	"github.com/tsavola/giobuffers/op"
	"github.com/tsavola/giobuffers/resource"
)

type ColorOp struct {
	Color color.NRGBA
}

func (x ColorOp) Create(b *flatbuffers.Builder) op.Op {
	flat.PaintColorOpStart(b)
	flat.PaintColorOpAddColor(b, flat.CreateColorNRGBA(b, x.Color.R, x.Color.G, x.Color.B, x.Color.A))
	offset := flat.PaintColorOpEnd(b)

	return op.Op{
		Type:   flat.OpPaintColor,
		Offset: offset,
	}
}

func (x ColorOp) Add(b *flatbuffers.Builder, o *op.Ops) {
	x.Create(b).Add(b, o)
}

type Image struct {
	resource.Resource
	typ    flat.Image
	decode []byte
	nrgba  *image.NRGBA
}

func NewImageDecode(c *resource.Context, id string, src []byte) *Image {
	return &Image{
		Resource: internal.MakeResource(c, id),
		typ:      flat.ImageImageDecode,
		decode:   src,
	}
}

func NewImageNRGBA(c *resource.Context, id string, src *image.NRGBA) *Image {
	return &Image{
		Resource: internal.MakeResource(c, id),
		typ:      flat.ImageImageNRGBA,
		nrgba:    src,
	}
}

func (x *Image) create(b *flatbuffers.Builder) flatbuffers.UOffsetT {
	switch x.typ {
	case flat.ImageImageDecode:
		data := b.CreateByteVector(x.decode)

		flat.ImageDecodeStart(b)
		flat.ImageDecodeAddData(b, data)
		return flat.ImageDecodeEnd(b)

	case flat.ImageImageNRGBA:
		pix := b.CreateByteVector(x.nrgba.Pix)

		flat.ImageNRGBAStart(b)
		flat.ImageNRGBAAddPix(b, pix)
		if n := x.nrgba.Stride; n != 0 {
			flat.ImageNRGBAAddStride(b, int32(n))
		}
		flat.ImageNRGBAAddRect(b, flat.CreateImageRectangle(b,
			int32(x.nrgba.Rect.Min.X), int32(x.nrgba.Rect.Min.Y),
			int32(x.nrgba.Rect.Max.X), int32(x.nrgba.Rect.Max.Y),
		))
		return flat.ImageNRGBAEnd(b)
	}

	panic(x.typ)
}

type ImageOp struct {
	src *Image
}

func NewImageOp(src *Image) ImageOp {
	return ImageOp{src}
}

func (x ImageOp) Create(b *flatbuffers.Builder) op.Op {
	idOff := b.CreateSharedString(x.src.ID)

	var srcOff flatbuffers.UOffsetT
	if x.src.NeedSync() {
		srcOff = x.src.create(b)
	}

	flat.PaintImageOpStart(b)
	flat.PaintImageOpAddSrcId(b, idOff)
	if srcOff != 0 {
		flat.PaintImageOpAddSrcType(b, x.src.typ)
		flat.PaintImageOpAddSrc(b, srcOff)
	}
	offset := flat.PaintImageOpEnd(b)

	x.src.Synced()

	return op.Op{
		Type:   flat.OpPaintImage,
		Offset: offset,
	}
}

func (x ImageOp) Add(b *flatbuffers.Builder, o *op.Ops) {
	x.Create(b).Add(b, o)
}

type PaintOp struct {
}

func (x PaintOp) Create(b *flatbuffers.Builder) op.Op {
	flat.PaintOpStart(b)
	offset := flat.PaintOpEnd(b)

	return op.Op{
		Type:   flat.OpPaint,
		Offset: offset,
	}
}

func (x PaintOp) Add(b *flatbuffers.Builder, o *op.Ops) {
	x.Create(b).Add(b, o)
}
