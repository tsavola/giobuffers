// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package flat

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type WidgetLabelLayout struct {
	_tab flatbuffers.Table
}

func GetRootAsWidgetLabelLayout(buf []byte, offset flatbuffers.UOffsetT) *WidgetLabelLayout {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &WidgetLabelLayout{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *WidgetLabelLayout) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *WidgetLabelLayout) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *WidgetLabelLayout) Label(obj *WidgetLabel) *WidgetLabel {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(WidgetLabel)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func (rcv *WidgetLabelLayout) Font(obj *TextFont) *TextFont {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(TextFont)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func (rcv *WidgetLabelLayout) Size(obj *UnitValue) *UnitValue {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		x := o + rcv._tab.Pos
		if obj == nil {
			obj = new(UnitValue)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func (rcv *WidgetLabelLayout) Text() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func WidgetLabelLayoutStart(builder *flatbuffers.Builder) {
	builder.StartObject(4)
}
func WidgetLabelLayoutAddLabel(builder *flatbuffers.Builder, label flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(label), 0)
}
func WidgetLabelLayoutAddFont(builder *flatbuffers.Builder, font flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(font), 0)
}
func WidgetLabelLayoutAddSize(builder *flatbuffers.Builder, size flatbuffers.UOffsetT) {
	builder.PrependStructSlot(2, flatbuffers.UOffsetT(size), 0)
}
func WidgetLabelLayoutAddText(builder *flatbuffers.Builder, text flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(3, flatbuffers.UOffsetT(text), 0)
}
func WidgetLabelLayoutEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}