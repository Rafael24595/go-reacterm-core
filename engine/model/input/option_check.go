package input

import (
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

type CheckOption struct {
	Id        string
	Status    bool
	Label     text.Frag
	Timestamp int64
}

func NewCheckOption(id string, option text.Frag) CheckOption {
	return CheckOption{
		Id:    id,
		Label: option,
	}
}

func FragsFromCheckOption(options ...CheckOption) []text.Frag {
	lines := make([]text.Frag, len(options))
	for i := range options {
		lines[i] = options[i].Label
	}
	return lines
}
