// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package flat

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type LabelLayout struct {
	_tab flatbuffers.Table
}

func GetRootAsLabelLayout(buf []byte, offset flatbuffers.UOffsetT) *LabelLayout {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &LabelLayout{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *LabelLayout) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *LabelLayout) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *LabelLayout) Label(obj *Label) *Label {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(Label)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func (rcv *LabelLayout) Constraints(obj *Constraints) *Constraints {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		x := o + rcv._tab.Pos
		if obj == nil {
			obj = new(Constraints)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func LabelLayoutStart(builder *flatbuffers.Builder) {
	builder.StartObject(2)
}
func LabelLayoutAddLabel(builder *flatbuffers.Builder, label flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(label), 0)
}
func LabelLayoutAddConstraints(builder *flatbuffers.Builder, constraints flatbuffers.UOffsetT) {
	builder.PrependStructSlot(1, flatbuffers.UOffsetT(constraints), 0)
}
func LabelLayoutEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}