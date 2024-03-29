// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package flat

import "strconv"

type WidgetFit byte

const (
	WidgetFitUnscaled  WidgetFit = 0
	WidgetFitContain   WidgetFit = 1
	WidgetFitCover     WidgetFit = 2
	WidgetFitScaleDown WidgetFit = 3
	WidgetFitFill      WidgetFit = 4
)

var EnumNamesWidgetFit = map[WidgetFit]string{
	WidgetFitUnscaled:  "Unscaled",
	WidgetFitContain:   "Contain",
	WidgetFitCover:     "Cover",
	WidgetFitScaleDown: "ScaleDown",
	WidgetFitFill:      "Fill",
}

var EnumValuesWidgetFit = map[string]WidgetFit{
	"Unscaled":  WidgetFitUnscaled,
	"Contain":   WidgetFitContain,
	"Cover":     WidgetFitCover,
	"ScaleDown": WidgetFitScaleDown,
	"Fill":      WidgetFitFill,
}

func (v WidgetFit) String() string {
	if s, ok := EnumNamesWidgetFit[v]; ok {
		return s
	}
	return "WidgetFit(" + strconv.FormatInt(int64(v), 10) + ")"
}
