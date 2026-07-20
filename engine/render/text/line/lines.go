package line

import (
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
)

func empty(
	order uint16,
	spec spec.Spec,
	size ...int,
) *Line {
	bufferSize := 0
	if len(size) > 0 {
		bufferSize = size[0]
	}

	return newLine(
		order,
		spec,
		make([]frag.Frag, 0, bufferSize),
	)
}

func Empty(size ...int) *Line {
	return empty(
		0,
		spec.Empty(),
		size...,
	)
}

func OrderedText(order uint16, text ...string) Line {
	return *newLine(
		order,
		spec.Empty(),
		frag.FromStrings(text...),
	)
}

func FromSpec(spec spec.Spec) Line {
	return *newLine(
		0,
		spec,
		make([]frag.Frag, 0),
	)
}

func FromMeta(other *Line, size ...int) *Line {
	return empty(
		other.Order,
		other.Spec,
		size...,
	)
}

func FromString(text ...string) Line {
	return *newLine(
		0,
		spec.Empty(),
		frag.FromStrings(text...),
	)
}

func FromFrags(frags ...frag.Frag) *Line {
	return newLine(
		0,
		spec.Empty(),
		frags,
	)
}
