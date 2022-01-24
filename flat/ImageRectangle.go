// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package flat

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type ImageRectangle struct {
	_tab flatbuffers.Struct
}

func (rcv *ImageRectangle) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *ImageRectangle) Table() flatbuffers.Table {
	return rcv._tab.Table
}

func (rcv *ImageRectangle) Min(obj *ImagePoint) *ImagePoint {
	if obj == nil {
		obj = new(ImagePoint)
	}
	obj.Init(rcv._tab.Bytes, rcv._tab.Pos+0)
	return obj
}
func (rcv *ImageRectangle) Max(obj *ImagePoint) *ImagePoint {
	if obj == nil {
		obj = new(ImagePoint)
	}
	obj.Init(rcv._tab.Bytes, rcv._tab.Pos+8)
	return obj
}

func CreateImageRectangle(builder *flatbuffers.Builder, min_x int32, min_y int32, max_x int32, max_y int32) flatbuffers.UOffsetT {
	builder.Prep(4, 16)
	builder.Prep(4, 8)
	builder.PrependInt32(max_y)
	builder.PrependInt32(max_x)
	builder.Prep(4, 8)
	builder.PrependInt32(min_y)
	builder.PrependInt32(min_x)
	return builder.Offset()
}
