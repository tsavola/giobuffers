// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package flat

import "strconv"

type Op byte

const (
	OpNONE        Op = 0
	OpMacro       Op = 1
	OpPaintColor  Op = 2
	OpPaint       Op = 3
	OpWidgetLabel Op = 4
)

var EnumNamesOp = map[Op]string{
	OpNONE:        "NONE",
	OpMacro:       "Macro",
	OpPaintColor:  "PaintColor",
	OpPaint:       "Paint",
	OpWidgetLabel: "WidgetLabel",
}

var EnumValuesOp = map[string]Op{
	"NONE":        OpNONE,
	"Macro":       OpMacro,
	"PaintColor":  OpPaintColor,
	"Paint":       OpPaint,
	"WidgetLabel": OpWidgetLabel,
}

func (v Op) String() string {
	if s, ok := EnumNamesOp[v]; ok {
		return s
	}
	return "Op(" + strconv.FormatInt(int64(v), 10) + ")"
}
