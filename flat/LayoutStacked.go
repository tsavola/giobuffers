// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package flat

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type LayoutStacked struct {
	_tab flatbuffers.Table
}

func GetRootAsLayoutStacked(buf []byte, offset flatbuffers.UOffsetT) *LayoutStacked {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &LayoutStacked{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *LayoutStacked) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *LayoutStacked) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *LayoutStacked) Widget(obj *OpNode) *OpNode {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(OpNode)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func LayoutStackedStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func LayoutStackedAddWidget(builder *flatbuffers.Builder, widget flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(widget), 0)
}
func LayoutStackedEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
