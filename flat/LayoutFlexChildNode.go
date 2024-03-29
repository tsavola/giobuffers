// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package flat

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type LayoutFlexChildNode struct {
	_tab flatbuffers.Table
}

func GetRootAsLayoutFlexChildNode(buf []byte, offset flatbuffers.UOffsetT) *LayoutFlexChildNode {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &LayoutFlexChildNode{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *LayoutFlexChildNode) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *LayoutFlexChildNode) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *LayoutFlexChildNode) ChildType() LayoutFlexChild {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return LayoutFlexChild(rcv._tab.GetByte(o + rcv._tab.Pos))
	}
	return 0
}

func (rcv *LayoutFlexChildNode) MutateChildType(n LayoutFlexChild) bool {
	return rcv._tab.MutateByteSlot(4, byte(n))
}

func (rcv *LayoutFlexChildNode) Child(obj *flatbuffers.Table) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		rcv._tab.Union(obj, o)
		return true
	}
	return false
}

func (rcv *LayoutFlexChildNode) Next(obj *LayoutFlexChildNode) *LayoutFlexChildNode {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(LayoutFlexChildNode)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func LayoutFlexChildNodeStart(builder *flatbuffers.Builder) {
	builder.StartObject(3)
}
func LayoutFlexChildNodeAddChildType(builder *flatbuffers.Builder, childType LayoutFlexChild) {
	builder.PrependByteSlot(0, byte(childType), 0)
}
func LayoutFlexChildNodeAddChild(builder *flatbuffers.Builder, child flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(child), 0)
}
func LayoutFlexChildNodeAddNext(builder *flatbuffers.Builder, next flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(next), 0)
}
func LayoutFlexChildNodeEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
