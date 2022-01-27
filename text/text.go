// Copyright (c) 2022 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package text

import (
	"github.com/tsavola/giobuffers/flat"
)

type Alignment = flat.TextAlignment

const (
	Start  = flat.TextAlignmentStart
	End    = flat.TextAlignmentEnd
	Middle = flat.TextAlignmentMiddle
)

type Font struct {
	Typeface string
	Variant  string
	Style    Style
	Weight   Weight
}

type Style = flat.TextStyle

const (
	Regular = flat.TextStyleRegular
	Italic  = flat.TextStyleItalic
)

type Weight int32

const (
	Thin       Weight = 100 - 400
	Hairline   Weight = Thin
	ExtraLight Weight = 200 - 400
	UltraLight Weight = ExtraLight
	Light      Weight = 300 - 400
	Normal     Weight = 400 - 400
	Medium     Weight = 500 - 400
	SemiBold   Weight = 600 - 400
	DemiBold   Weight = SemiBold
	Bold       Weight = 700 - 400
	ExtraBold  Weight = 800 - 400
	UltraBold  Weight = ExtraBold
	Black      Weight = 900 - 400
	Heavy      Weight = Black
	ExtraBlack Weight = 950 - 400
	UltraBlack Weight = ExtraBlack
)
