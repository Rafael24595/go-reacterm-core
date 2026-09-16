package event

import (
	"github.com/Rafael24595/go-reacterm-core/engine/helper/runes"
	"github.com/Rafael24595/go-reacterm-core/engine/model/offset"
)

type ActionKind int

const (
	Insert ActionKind = iota

	DeleteBackward
	DeleteForward

	Cut
	Paste
)

type textAction struct {
	kind      ActionKind
	start     offset.Offset
	end       offset.Offset
	delete    string
	insert    string
	timestamp int64
}

type mergeAction struct {
	kind   ActionKind
	origin offset.Offset
	extent offset.Offset
	probe  offset.Offset
	delete []string
	insert []string
}

func (m *mergeAction) len() offset.Offset {
	var n offset.Offset
	for _, t := range m.insert {
		n += runes.MeasureOffset(t)
	}
	return n
}
