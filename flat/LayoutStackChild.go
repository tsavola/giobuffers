// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package flat

import "strconv"

type LayoutStackChild byte

const (
	LayoutStackChildNONE     LayoutStackChild = 0
	LayoutStackChildExpanded LayoutStackChild = 1
	LayoutStackChildStacked  LayoutStackChild = 2
)

var EnumNamesLayoutStackChild = map[LayoutStackChild]string{
	LayoutStackChildNONE:     "NONE",
	LayoutStackChildExpanded: "Expanded",
	LayoutStackChildStacked:  "Stacked",
}

var EnumValuesLayoutStackChild = map[string]LayoutStackChild{
	"NONE":     LayoutStackChildNONE,
	"Expanded": LayoutStackChildExpanded,
	"Stacked":  LayoutStackChildStacked,
}

func (v LayoutStackChild) String() string {
	if s, ok := EnumNamesLayoutStackChild[v]; ok {
		return s
	}
	return "LayoutStackChild(" + strconv.FormatInt(int64(v), 10) + ")"
}