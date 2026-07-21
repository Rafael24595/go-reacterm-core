package line

import (
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
)

func Empty(size ...int) Line {
	bufferSize := 0
	if len(size) > 0 {
		bufferSize = size[0]
	}

	return New(
		0,
		spec.Empty(),
		make([]frag.Frag, 0, bufferSize),
	)
}

func TextSpec(text string, specs ...spec.Spec) Line {
	return New(
		0,
		spec.Merge(specs...),
		frag.FromStrings(text),
	)
}

func TextOrdered(order uint16, text ...string) Line {
	return New(
		order,
		spec.Empty(),
		frag.FromStrings(text...),
	)
}

func FromSpec(spec spec.Spec) Line {
	return New(
		0,
		spec,
		make([]frag.Frag, 0),
	)
}

func FromString(text ...string) Line {
	return New(
		0,
		spec.Empty(),
		frag.FromStrings(text...),
	)
}

func FromFrags(frags ...frag.Frag) Line {
	return New(
		0,
		spec.Empty(),
		frags,
	)
}
