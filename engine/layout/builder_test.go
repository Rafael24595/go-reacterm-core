package layout

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

func dummyComposer(
	uiState *state.UIState,
	size winsize.Winsize,
	vm viewmodel.ViewModel,
) (*state.UIState, []line.Line) {
	return uiState, nil
}

func TestLayoutBuilder_Basic(t *testing.T) {
	builder := NewBuilder(dummyComposer)
	l := builder.ToLayout()

	assert.False(t, l.Compose == nil)

	ws := winsize.New(10, 20)
	vm := viewmodel.New()

	state, lines := l.Compose(nil, ws, *vm)

	assert.Nil(t, state)
	assert.Nil(t, lines)
}

func TestLayoutBuilder_WithTransformer(t *testing.T) {
	var resultWs winsize.Winsize

	composer := func(
		uiState *state.UIState,
		size winsize.Winsize,
		vm viewmodel.ViewModel,
	) (*state.UIState, []line.Line) {
		resultWs = size
		return uiState, nil
	}

	doubleSizeTransformer := func(size winsize.Winsize) winsize.Winsize {
		return winsize.Winsize{
			Rows: size.Rows * 2,
			Cols: size.Cols * 2,
		}
	}

	layout := NewBuilder(composer).
		Transformer(doubleSizeTransformer).
		ToLayout()

	ws := winsize.New(10, 20)
	vm := viewmodel.New()

	layout.Compose(nil, ws, *vm)

	wantWs := winsize.New(20, 40)

	assert.Equal(t, wantWs, resultWs)
}
