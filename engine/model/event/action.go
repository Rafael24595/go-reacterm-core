package event

import (
	"github.com/Rafael24595/go-reacterm-core/engine/helper/runes"
	"github.com/Rafael24595/go-reacterm-core/engine/model/offset"
)

// ActionKind defines the nature of a text mutation.
type ActionKind int

const (
	// Insert represents a text insertion operation.
	Insert ActionKind = iota

	// DeleteBackward represents a backward deletion operation, which is a special case of deletion.
	DeleteBackward
	// DeleteForward represents a forward deletion operation, which is a special case of deletion.
	DeleteForward

	// Cut represents a cut operation, which is a special case of deletion.
	Cut
	// Paste represents a paste operation, which is a special case of insertion.
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
