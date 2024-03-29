// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package flat

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type PaintImageOp struct {
	_tab flatbuffers.Table
}

func GetRootAsPaintImageOp(buf []byte, offset flatbuffers.UOffsetT) *PaintImageOp {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &PaintImageOp{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *PaintImageOp) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *PaintImageOp) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *PaintImageOp) SrcId() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *PaintImageOp) SrcType() Image {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return Image(rcv._tab.GetByte(o + rcv._tab.Pos))
	}
	return 0
}

func (rcv *PaintImageOp) MutateSrcType(n Image) bool {
	return rcv._tab.MutateByteSlot(6, byte(n))
}

func (rcv *PaintImageOp) Src(obj *flatbuffers.Table) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		rcv._tab.Union(obj, o)
		return true
	}
	return false
}

func PaintImageOpStart(builder *flatbuffers.Builder) {
	builder.StartObject(3)
}
func PaintImageOpAddSrcId(builder *flatbuffers.Builder, srcId flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(srcId), 0)
}
func PaintImageOpAddSrcType(builder *flatbuffers.Builder, srcType Image) {
	builder.PrependByteSlot(1, byte(srcType), 0)
}
func PaintImageOpAddSrc(builder *flatbuffers.Builder, src flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(src), 0)
}
func PaintImageOpEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
