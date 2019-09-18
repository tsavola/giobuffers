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

	regular := b.CreateString("Go Regular")

	flat.FaceStart(b)
	flat.FaceAddFont(b, regular)
	flat.FaceAddSize(b, flat.CreateValue(b, 72, flat.UnitSp))
	face := flat.FaceEnd(b)

	message := b.CreateString("Hello, Gio")

	flat.ColorOpStart(b)
	flat.ColorOpAddColor(b, flat.CreateColorRGBA(b, 127, 0, 0, 255))
	maroon := flat.ColorOpEnd(b)

	flat.OpNodeStart(b)
	flat.OpNodeAddOpType(b, flat.OpColor)
	flat.OpNodeAddOp(b, maroon)
	material := flat.OpNodeEnd(b)

	flat.LabelStart(b)
	flat.LabelAddMaterial(b, material)
	flat.LabelAddFace(b, face)
	flat.LabelAddAlignment(b, flat.AlignmentMiddle)
	flat.LabelAddText(b, message)
	label := flat.LabelEnd(b)

	flat.LabelLayoutStart(b)
	flat.LabelLayoutAddLabel(b, label)
	flat.LabelLayoutAddConstraints(b, flat.CreateConstraints(b, 0, 640, 0, 480))
	layout := flat.LabelLayoutEnd(b)

	flat.OpNodeStart(b)
	flat.OpNodeAddOpType(b, flat.OpLabelLayout)
	flat.OpNodeAddOp(b, layout)
	ops := flat.OpNodeEnd(b)

	flat.OpNodeStart(b)
	flat.OpNodeAddNext(b, ops)
	flat.OpNodeAddOpType(b, flat.OpMacro)
	flat.OpNodeAddOp(b, material)
	ops = flat.OpNodeEnd(b)

	b.Finish(ops)
	return b.FinishedBytes()
}
