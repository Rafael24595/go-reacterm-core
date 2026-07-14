package input

import "github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"

type CheckOption struct {
	Id        string
	Status    bool
	Label     frag.Frag
	Timestamp int64
}

func NewCheckOption(id string, option frag.Frag) CheckOption {
	return CheckOption{
		Id:    id,
		Label: option,
	}
}

func FragsFromCheckOption(options ...CheckOption) []frag.Frag {
	lines := make([]frag.Frag, len(options))
	for i := range options {
		lines[i] = options[i].Label
	}
	return lines
}
