// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package flat

import "strconv"

type LayoutAlignment byte

const (
	LayoutAlignmentStart    LayoutAlignment = 0
	LayoutAlignmentEnd      LayoutAlignment = 1
	LayoutAlignmentMiddle   LayoutAlignment = 2
	LayoutAlignmentBaseline LayoutAlignment = 3
)

var EnumNamesLayoutAlignment = map[LayoutAlignment]string{
	LayoutAlignmentStart:    "Start",
	LayoutAlignmentEnd:      "End",
	LayoutAlignmentMiddle:   "Middle",
	LayoutAlignmentBaseline: "Baseline",
}

var EnumValuesLayoutAlignment = map[string]LayoutAlignment{
	"Start":    LayoutAlignmentStart,
	"End":      LayoutAlignmentEnd,
	"Middle":   LayoutAlignmentMiddle,
	"Baseline": LayoutAlignmentBaseline,
}

func (v LayoutAlignment) String() string {
	if s, ok := EnumNamesLayoutAlignment[v]; ok {
		return s
	}
	return "LayoutAlignment(" + strconv.FormatInt(int64(v), 10) + ")"
}
