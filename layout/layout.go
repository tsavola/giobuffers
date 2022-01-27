// Copyright (c) 2022 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package layout

import (
	"gioui.org/unit"
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/tsavola/giobuffers/flat"
	"github.com/tsavola/giobuffers/op"
)

type Alignment = flat.LayoutAlignment

const (
	Start    = flat.LayoutAlignmentStart
	End      = flat.LayoutAlignmentEnd
	Middle   = flat.LayoutAlignmentMiddle
	Baseline = flat.LayoutAlignmentBaseline
)

type Axis = flat.LayoutAxis

const (
	Horizontal = flat.LayoutAxisHorizontal
	Vertical   = flat.LayoutAxisVertical
)

type Spacing = flat.LayoutSpacing

const (
	SpaceEnd     = flat.LayoutSpacingEnd
	SpaceStart   = flat.LayoutSpacingStart
	SpaceSides   = flat.LayoutSpacingSides
	SpaceAround  = flat.LayoutSpacingAround
	SpaceBetween = flat.LayoutSpacingBetween
	SpaceEvenly  = flat.LayoutSpacingEvenly
)

type Direction flat.LayoutDirection

const (
	NW     = Direction(flat.LayoutDirectionNW)
	N      = Direction(flat.LayoutDirectionN)
	NE     = Direction(flat.LayoutDirectionNE)
	E      = Direction(flat.LayoutDirectionE)
	SE     = Direction(flat.LayoutDirectionSE)
	S      = Direction(flat.LayoutDirectionS)
	SW     = Direction(flat.LayoutDirectionSW)
	W      = Direction(flat.LayoutDirectionW)
	Center = Direction(flat.LayoutDirectionCenter)
)

func (x Direction) CreateLayout(b *flatbuffers.Builder, widget op.Ops) op.Op {
	flat.LayoutDirectionLayoutStart(b)
	flat.LayoutDirectionLayoutAddDirection(b, flat.LayoutDirection(x))
	flat.LayoutDirectionLayoutAddWidget(b, widget.LastNode)
	offset := flat.LayoutDirectionLayoutEnd(b)

	return op.Op{
		Type:   flat.OpLayoutDirection,
		Offset: offset,
	}
}

func (x Direction) AddLayout(b *flatbuffers.Builder, o *op.Ops, widget op.Ops) op.Ops {
	return x.CreateLayout(b, widget).Add(b, o)
}

type Flex struct {
	Axis      Axis
	Spacing   Spacing
	Alignment Alignment
	WeightSum float32
}

func (x Flex) create(b *flatbuffers.Builder) flatbuffers.UOffsetT {
	flat.LayoutFlexStart(b)
	if x.Axis != 0 {
		flat.LayoutFlexAddAxis(b, x.Axis)
	}
	if x.Spacing != 0 {
		flat.LayoutFlexAddSpacing(b, x.Spacing)
	}
	if x.Alignment != 0 {
		flat.LayoutFlexAddAlignment(b, x.Alignment)
	}
	if x.WeightSum != 0 {
		flat.LayoutFlexAddWeightSum(b, x.WeightSum)
	}
	return flat.LayoutFlexEnd(b)
}

func (x Flex) CreateLayout(b *flatbuffers.Builder, children ...FlexChild) op.Op {
	flexOff := x.create(b)

	var childOps op.Ops
	for i := len(children) - 1; i >= 0; i-- {
		children[i].addPrevious(b, &childOps)
	}

	flat.LayoutFlexLayoutStart(b)
	flat.LayoutFlexLayoutAddFlex(b, flexOff)
	flat.LayoutFlexLayoutAddChildren(b, childOps.LastNode)
	offset := flat.LayoutFlexLayoutEnd(b)

	return op.Op{
		Type:   flat.OpLayoutFlex,
		Offset: offset,
	}
}

func (x Flex) AddLayout(b *flatbuffers.Builder, o *op.Ops, children ...FlexChild) op.Ops {
	return x.CreateLayout(b, children...).Add(b, o)
}

type FlexChild struct {
	typ    flat.LayoutFlexChild
	weight float32
	widget op.Ops
}

func Flexed(weight float32, widget op.Ops) FlexChild {
	return FlexChild{
		typ:    flat.LayoutFlexChildFlexed,
		weight: weight,
		widget: widget,
	}
}

func Rigid(widget op.Ops) FlexChild {
	return FlexChild{
		typ:    flat.LayoutFlexChildRigid,
		widget: widget,
	}
}

func (c FlexChild) create(b *flatbuffers.Builder) flatbuffers.UOffsetT {
	switch c.typ {
	case flat.LayoutFlexChildFlexed:
		flat.LayoutFlexedStart(b)
		flat.LayoutFlexedAddWeight(b, c.weight)
		flat.LayoutFlexedAddWidget(b, c.widget.LastNode)
		return flat.LayoutFlexedEnd(b)

	case flat.LayoutFlexChildRigid:
		flat.LayoutRigidStart(b)
		flat.LayoutRigidAddWidget(b, c.widget.LastNode)
		return flat.LayoutRigidEnd(b)
	}

	panic(c.typ)
}

func (c FlexChild) addPrevious(b *flatbuffers.Builder, children *op.Ops) {
	offset := c.create(b)

	flat.LayoutFlexChildNodeStart(b)
	flat.LayoutFlexChildNodeAddChildType(b, c.typ)
	flat.LayoutFlexChildNodeAddChild(b, offset)
	if children.LastNode != 0 {
		flat.LayoutFlexChildNodeAddNext(b, children.LastNode)
	}
	children.LastNode = flat.LayoutFlexChildNodeEnd(b)
}

type Inset struct {
	Top    unit.Value
	Bottom unit.Value
	Left   unit.Value
	Right  unit.Value
}

func UniformInset(v unit.Value) Inset {
	return Inset{v, v, v, v}
}

func (x Inset) create(b *flatbuffers.Builder) flatbuffers.UOffsetT {
	flat.LayoutInsetStart(b)
	flat.LayoutInsetAddTop(b, flat.CreateUnitValue(b, x.Top.V, flat.Unit(x.Top.U)))
	flat.LayoutInsetAddBottom(b, flat.CreateUnitValue(b, x.Bottom.V, flat.Unit(x.Bottom.U)))
	flat.LayoutInsetAddLeft(b, flat.CreateUnitValue(b, x.Left.V, flat.Unit(x.Left.U)))
	flat.LayoutInsetAddRight(b, flat.CreateUnitValue(b, x.Right.V, flat.Unit(x.Right.U)))
	return flat.LayoutInsetEnd(b)
}

func (x Inset) CreateLayout(b *flatbuffers.Builder, widget op.Ops) op.Op {
	insetOff := x.create(b)

	flat.LayoutInsetLayoutStart(b)
	flat.LayoutInsetLayoutAddInset(b, insetOff)
	flat.LayoutInsetLayoutAddWidget(b, widget.LastNode)
	offset := flat.LayoutInsetLayoutEnd(b)

	return op.Op{
		Type:   flat.OpLayoutInset,
		Offset: offset,
	}
}

func (x Inset) AddLayout(b *flatbuffers.Builder, o *op.Ops, widget op.Ops) op.Ops {
	return x.CreateLayout(b, widget).Add(b, o)
}

type Spacer struct {
	Width  unit.Value
	Height unit.Value
}

func (x Spacer) create(b *flatbuffers.Builder) flatbuffers.UOffsetT {
	flat.LayoutSpacerStart(b)
	flat.LayoutSpacerAddWidth(b, flat.CreateUnitValue(b, x.Width.V, flat.Unit(x.Width.U)))
	flat.LayoutSpacerAddHeight(b, flat.CreateUnitValue(b, x.Height.V, flat.Unit(x.Height.U)))
	return flat.LayoutSpacerEnd(b)
}

func (x Spacer) CreateLayout(b *flatbuffers.Builder) op.Op {
	spacerOff := x.create(b)

	flat.LayoutSpacerLayoutStart(b)
	flat.LayoutSpacerLayoutAddSpacer(b, spacerOff)
	offset := flat.LayoutSpacerLayoutEnd(b)

	return op.Op{
		Type:   flat.OpLayoutSpacer,
		Offset: offset,
	}
}

func (x Spacer) AddLayout(b *flatbuffers.Builder, o *op.Ops) op.Ops {
	return x.CreateLayout(b).Add(b, o)
}

type Stack struct {
	Alignment Direction
}

func (x Stack) create(b *flatbuffers.Builder) flatbuffers.UOffsetT {
	flat.LayoutStackStart(b)
	if x.Alignment != 0 {
		flat.LayoutStackAddAlignment(b, flat.LayoutDirection(x.Alignment))
	}
	return flat.LayoutStackEnd(b)
}

func (x Stack) CreateLayout(b *flatbuffers.Builder, children ...StackChild) op.Op {
	stackOff := x.create(b)

	var childOps op.Ops
	for i := len(children) - 1; i >= 0; i-- {
		children[i].addPrevious(b, &childOps)
	}

	flat.LayoutStackLayoutStart(b)
	flat.LayoutStackLayoutAddStack(b, stackOff)
	flat.LayoutStackLayoutAddChildren(b, childOps.LastNode)
	offset := flat.LayoutStackLayoutEnd(b)

	return op.Op{
		Type:   flat.OpLayoutStack,
		Offset: offset,
	}
}

func (x Stack) AddLayout(b *flatbuffers.Builder, o *op.Ops, children ...StackChild) op.Ops {
	return x.CreateLayout(b, children...).Add(b, o)
}

type StackChild struct {
	typ    flat.LayoutStackChild
	widget op.Ops
}

func Expanded(widget op.Ops) StackChild {
	return StackChild{
		typ:    flat.LayoutStackChildExpanded,
		widget: widget,
	}
}

func Stacked(widget op.Ops) StackChild {
	return StackChild{
		typ:    flat.LayoutStackChildStacked,
		widget: widget,
	}
}

func (c StackChild) create(b *flatbuffers.Builder) flatbuffers.UOffsetT {
	switch c.typ {
	case flat.LayoutStackChildExpanded:
		flat.LayoutExpandedStart(b)
		flat.LayoutExpandedAddWidget(b, c.widget.LastNode)
		return flat.LayoutExpandedEnd(b)

	case flat.LayoutStackChildStacked:
		flat.LayoutStackedStart(b)
		flat.LayoutStackedAddWidget(b, c.widget.LastNode)
		return flat.LayoutStackedEnd(b)
	}

	panic(c.typ)
}

func (c StackChild) addPrevious(b *flatbuffers.Builder, children *op.Ops) {
	offset := c.create(b)

	flat.LayoutStackChildNodeStart(b)
	flat.LayoutStackChildNodeAddChildType(b, c.typ)
	flat.LayoutStackChildNodeAddChild(b, offset)
	if children.LastNode != 0 {
		flat.LayoutStackChildNodeAddNext(b, children.LastNode)
	}
	children.LastNode = flat.LayoutStackChildNodeEnd(b)
}
