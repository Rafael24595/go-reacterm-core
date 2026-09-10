package layout

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

// LayoutBuilder constructs a Layout instance with optional window size transformations.
type LayoutBuilder struct {
	transformer winsize.Transformer
	compose     Composer
}

// NewBuilder initializes a LayoutBuilder with a non-nil Composer.
func NewBuilder(composer Composer) *LayoutBuilder {
	assert.False(composer == nil, "Composer function cannot be nil")

	return &LayoutBuilder{
		compose: composer,
	}
}

// Transformer attaches a winsize.Transformer to pre-process window dimensions before layout composition.
func (b *LayoutBuilder) Transformer(transformer winsize.Transformer) *LayoutBuilder {
	b.transformer = transformer
	return b
}

// ToLayout builds and returns the final Layout struct.
func (b *LayoutBuilder) ToLayout() Layout {
	compose := b.compose
	if b.transformer != nil {
		compose = wrapTransformer(compose, b.transformer)
	}

	return Layout{
		Compose: compose,
	}
}

func wrapTransformer(compose Composer, transformer winsize.Transformer) Composer {
	return func(uiState *state.UIState, size winsize.Winsize, vm viewmodel.ViewModel) (*state.UIState, []line.Line) {
		newSize := transformer(size)
		return compose(uiState, newSize, vm)
	}
}
