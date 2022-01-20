// Copyright (c) 2019 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package giobuffers

import (
	"testing"

	"gioui.org/layout"
	"gioui.org/op"
	"github.com/tsavola/giobuffers/internal/hello"
)

var data = hello.Marshal()

func BenchmarkUnmarshal(b *testing.B) {
	b.SetBytes(int64(len(data)))

	var ops op.Ops

	var u Unmarshaler

	for i := 0; i < b.N; i++ {
		ops.Reset()
		gtx := layout.Context{Ops: &ops}

		if err := u.Unmarshal(gtx, data); err != nil {
			b.Fatal(err)
		}
	}
}
