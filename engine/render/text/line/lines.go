package line

import (
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
)

func Empty(size ...int) *Line {
	bufferSize := 0
	if len(size) > 0 {
		bufferSize = size[0]
	}

	return FromFrags(
		make([]frag.Frag, 0, bufferSize)...,
	)
}

func FromMeta(other *Line, size ...int) *Line {
	return Empty(size...).
		CopyMeta(other)
}

func FromFrags(frags ...frag.Frag) *Line {
	return &Line{
		Text: frags,
		Spec: spec.Empty(),
	}
}
