// Copyright (c) 2022 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package op

import (
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/tsavola/giobuffers/flat"
)

type Op struct {
	Type   flat.Op
	Offset flatbuffers.UOffsetT
}

func (x Op) Add(b *flatbuffers.Builder, o *Ops) Ops {
	if o == nil {
		o = new(Ops)
	}

	flat.OpNodeStart(b)
	flat.OpNodeAddOpType(b, x.Type)
	flat.OpNodeAddOp(b, x.Offset)
	if o.LastNode != 0 {
		flat.OpNodeAddPrevious(b, o.LastNode)
	}
	o.LastNode = flat.OpNodeEnd(b)

	return *o
}

type Ops struct {
	LastNode flatbuffers.UOffsetT
}

func (o *Ops) Reset() {
	*o = Ops{}
}
