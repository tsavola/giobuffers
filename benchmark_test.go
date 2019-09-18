// Copyright (c) 2019 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package giobuffers

import (
	"testing"
	"time"

	"gioui.org/ui"
	"gioui.org/ui/measure"
	"github.com/tsavola/giobuffers/internal/hello"
)

var data = hello.Marshal()

type config struct{}

func (config) Now() time.Time    { return time.Now() }
func (config) Px(v ui.Value) int { return int(v.V) }

func BenchmarkUnmarshal(b *testing.B) {
	b.SetBytes(int64(len(data)))

	var faces measure.Faces
	faces.Reset(config{})

	var ops ui.Ops

	var u Unmarshaler

	for i := 0; i < b.N; i++ {
		ops.Reset()

		if err := u.Unmarshal(data, &ops, &faces); err != nil {
			b.Fatal(err)
		}
	}
}
