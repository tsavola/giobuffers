// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package flat

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type ColorOp struct {
	_tab flatbuffers.Table
}

func GetRootAsColorOp(buf []byte, offset flatbuffers.UOffsetT) *ColorOp {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &ColorOp{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *ColorOp) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *ColorOp) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *ColorOp) Color(obj *ColorRGBA) *ColorRGBA {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		x := o + rcv._tab.Pos
		if obj == nil {
			obj = new(ColorRGBA)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func ColorOpStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func ColorOpAddColor(builder *flatbuffers.Builder, color flatbuffers.UOffsetT) {
	builder.PrependStructSlot(0, flatbuffers.UOffsetT(color), 0)
}
func ColorOpEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
