// Copyright (c) 2019 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"github.com/tsavola/giobuffers"
	"github.com/tsavola/giobuffers/internal/hello"
)

func main() {
	data := hello.Marshal()
	go func() {
		w := app.NewWindow()
		err := run(w, data)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(w *app.Window, data []byte) error {
	var u giobuffers.Unmarshaler
	var ops op.Ops
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			if _, err := u.Unmarshal(gtx, data); err != nil {
				return err
			}
			e.Frame(gtx.Ops)
		}
	}
}
