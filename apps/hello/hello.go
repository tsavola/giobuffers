// Copyright (c) 2019 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"

	"gioui.org/ui"
	"gioui.org/ui/app"
	"gioui.org/ui/measure"
	"github.com/tsavola/giobuffers"
	"github.com/tsavola/giobuffers/internal/hello"
)

func main() {
	data := hello.Marshal()
	go func() {
		w := app.NewWindow()
		if err := loop(w, data); err != nil {
			log.Fatal(err)
		}
	}()
	app.Main()
}

func loop(w *app.Window, data []byte) error {
	var cfg app.Config
	var faces measure.Faces
	ops := new(ui.Ops)
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case app.DestroyEvent:
			return e.Err
		case app.UpdateEvent:
			cfg = e.Config
			faces.Reset(&cfg)
			ops.Reset()
			if err := giobuffers.Unmarshal(data, ops, &faces); err != nil {
				return err
			}
			w.Update(ops)
		}
	}
}
