package pagination

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/pager/predicate"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"

	screen_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/screen"
	text_test "github.com/Rafael24595/go-reacterm-core/test/engine/render/text"
)

func TestPagination_ToNode(t *testing.T) {
	name := "base"
	mock := screen_test.MockByName(name)

	node := New(mock).ToNode()
	screen_test.Helper_ToNode(t, node)

	assert.Equal(t, node.Name, name)
}

func TestPagination_Propagate(t *testing.T) {
	name := "base"
	mock := screen_test.MockByName(name)

	node := New(mock).ToNode()
	screen_test.Helper_Propagate(t, name, 0, node)
}

func TestPagination_LocalTick(t *testing.T) {
	uiState := state.NewUIState()
	base := screen_test.MockByName("base")
	page := New(base)
	node := page.ToNode()

	uiState.Pager.TargetPage = 0
	result := node.Screen.Tick(uiState, screen.Event{Key: *key.NewKeyCode(key.ActionArrowLeft)})
	assert.Equal(t, result.Pager.TargetPage, 0)

	result = node.Screen.Tick(uiState, screen.Event{Key: *key.NewKeyCode(key.ActionArrowRight)})
	assert.Equal(t, result.Pager.TargetPage, 1)
}

func TestPagination_ViewFooter(t *testing.T) {
	uiState := state.NewUIState()
	uiState.Pager.ActualPage = 3

	base := screen_test.MockNode{
		Name: "base",
		View: func(_ state.UIState) viewmodel.ViewModel {
			vm := viewmodel.New()
			vm.Pager.SetPredicate(predicate.Page())
			return *vm
		},
	}

	page := New(base.ToNode())
	vm := page.view(*uiState)

	footer := vm.Footer.ToUnit()
	footer.Drawable.Boot()

	lines, _ := footer.Drawable.Draw(winsize.Winsize{
		Rows: 3,
		Cols: 10,
	})

	assert.True(t, len(lines) > 0)
	assert.Inside(t, "page: 3", text_test.LineToString(lines[0]))
}

func TestPagination_TickDelegates(t *testing.T) {
	called := false

	base := screen_test.MockNode{
		Name: "base",
		Tick: func(s *state.UIState, e screen.Event) screen.Result {
			called = true
			return screen.EmptyResult()
		},
	}

	page := New(base.ToNode())
	node := page.ToNode()

	node.Screen.Tick(state.NewUIState(), screen.Event{Key: *key.NewKeyRune('x')})

	assert.True(t, called, "screen.Tick should be called")
}

func TestPagination_PageNeverNegative(t *testing.T) {
	uiState := state.NewUIState()
	uiState.Pager.TargetPage = 0

	base := screen_test.MockByName("base")
	page := New(base)
	node := page.ToNode()

	node.Screen.Tick(uiState, screen.Event{Key: *key.NewKeyCode(key.ActionArrowLeft)})
	assert.Equal(t, uiState.Pager.TargetPage, 0)
}
