// Copyright (c) 2019 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hello

import (
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/tsavola/giobuffers/flat"
)

func Marshal() []byte {
	b := flatbuffers.NewBuilder(0)

	var opList flatbuffers.UOffsetT
	addOp := func(opType flat.Op, bufOffset flatbuffers.UOffsetT) {
		flat.OpNodeStart(b)
		flat.OpNodeAddOpType(b, opType)
		flat.OpNodeAddOp(b, bufOffset)
		flat.OpNodeAddPrevious(b, opList)
		opList = flat.OpNodeEnd(b)
	}

	flat.PaintLinearGradientOpStart(b)
	flat.PaintLinearGradientOpAddStop1(b, flat.CreateF32Point(b, 0, 0))
	flat.PaintLinearGradientOpAddColor1(b, flat.CreateColorNRGBA(b, 16, 32, 48, 255))
	flat.PaintLinearGradientOpAddStop2(b, flat.CreateF32Point(b, 32, 256))
	flat.PaintLinearGradientOpAddColor2(b, flat.CreateColorNRGBA(b, 0, 0, 0, 255))
	addOp(flat.OpPaintLinearGradient, flat.PaintLinearGradientOpEnd(b))

	flat.PaintOpStart(b)
	addOp(flat.OpPaint, flat.PaintOpEnd(b))

	flat.PaintColorOpStart(b)
	flat.PaintColorOpAddColor(b, flat.CreateColorNRGBA(b, 200, 200, 200, 255))
	addOp(flat.OpPaintColor, flat.PaintColorOpEnd(b))

	flat.WidgetLabelStart(b)
	flat.WidgetLabelAddAlignment(b, flat.TextAlignmentMiddle)
	label := flat.WidgetLabelEnd(b)

	flat.TextFontStart(b)
	flat.TextFontAddStyle(b, flat.TextStyleItalic)
	flat.TextFontAddWeight(b, 550)
	font := flat.TextFontEnd(b)

	text := b.CreateString("Hello, Gio")

	flat.WidgetLabelLayoutStart(b)
	flat.WidgetLabelLayoutAddLabel(b, label)
	flat.WidgetLabelLayoutAddFont(b, font)
	flat.WidgetLabelLayoutAddSize(b, flat.CreateUnitValue(b, 16, flat.UnitDp))
	flat.WidgetLabelLayoutAddText(b, text)
	addOp(flat.OpWidgetLabel, flat.WidgetLabelLayoutEnd(b))

	b.Finish(opList)
	return b.FinishedBytes()
}
