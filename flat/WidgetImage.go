// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package flat

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type WidgetImage struct {
	_tab flatbuffers.Table
}

func GetRootAsWidgetImage(buf []byte, offset flatbuffers.UOffsetT) *WidgetImage {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &WidgetImage{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *WidgetImage) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *WidgetImage) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *WidgetImage) Src(obj *PaintImageOp) *PaintImageOp {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(PaintImageOp)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func (rcv *WidgetImage) Fit() WidgetFit {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return WidgetFit(rcv._tab.GetByte(o + rcv._tab.Pos))
	}
	return 0
}

func (rcv *WidgetImage) MutateFit(n WidgetFit) bool {
	return rcv._tab.MutateByteSlot(6, byte(n))
}

func (rcv *WidgetImage) Position() LayoutDirection {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return LayoutDirection(rcv._tab.GetByte(o + rcv._tab.Pos))
	}
	return 0
}

func (rcv *WidgetImage) MutatePosition(n LayoutDirection) bool {
	return rcv._tab.MutateByteSlot(8, byte(n))
}

func (rcv *WidgetImage) Scale() float32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetFloat32(o + rcv._tab.Pos)
	}
	return 0.0
}

func (rcv *WidgetImage) MutateScale(n float32) bool {
	return rcv._tab.MutateFloat32Slot(10, n)
}

func WidgetImageStart(builder *flatbuffers.Builder) {
	builder.StartObject(4)
}
func WidgetImageAddSrc(builder *flatbuffers.Builder, src flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(src), 0)
}
func WidgetImageAddFit(builder *flatbuffers.Builder, fit WidgetFit) {
	builder.PrependByteSlot(1, byte(fit), 0)
}
func WidgetImageAddPosition(builder *flatbuffers.Builder, position LayoutDirection) {
	builder.PrependByteSlot(2, byte(position), 0)
}
func WidgetImageAddScale(builder *flatbuffers.Builder, scale float32) {
	builder.PrependFloat32Slot(3, scale, 0.0)
}
func WidgetImageEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
