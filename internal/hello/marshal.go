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
		flat.OpNodeAddNext(b, opList)
		opList = flat.OpNodeEnd(b)
	}

	addColorOp := func(color flatbuffers.UOffsetT) {
		flat.PaintColorOpStart(b)
		flat.PaintColorOpAddColor(b, color)
		addOp(flat.OpPaintColor, flat.PaintColorOpEnd(b))
	}

	var f lazyReverse
	f.add(func() {
		flat.PaintLinearGradientOpStart(b)
		flat.PaintLinearGradientOpAddStop1(b, flat.CreateF32Point(b, 0, 0))
		flat.PaintLinearGradientOpAddColor1(b, flat.CreateColorNRGBA(b, 16, 32, 48, 255))
		flat.PaintLinearGradientOpAddStop2(b, flat.CreateF32Point(b, 32, 256))
		flat.PaintLinearGradientOpAddColor2(b, flat.CreateColorNRGBA(b, 0, 0, 0, 255))
		addOp(flat.OpPaintLinearGradient, flat.PaintLinearGradientOpEnd(b))
	})
	f.add(func() {
		flat.PaintOpStart(b)
		addOp(flat.OpPaint, flat.PaintOpEnd(b))
	})
	f.add(func() {
		addColorOp(flat.CreateColorNRGBA(b, 200, 200, 200, 255))
	})
	f.add(func() {
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
	})
	f.eval()

	b.Finish(opList)
	return b.FinishedBytes()
}

type lazyReverse struct {
	first *funcNode
}

func (lr *lazyReverse) add(f func()) {
	lr.first = &funcNode{f, lr.first}
}

func (lr *lazyReverse) eval() {
	lr.first.eval()
}

type funcNode struct {
	f    func()
	next *funcNode
}

func (fn *funcNode) eval() {
	if fn != nil {
		fn.f()
		fn.next.eval()
	}
}
