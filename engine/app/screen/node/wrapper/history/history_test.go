package history

import (
	"fmt"
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"

	screen_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/screen"
)

func TestHistory_ToNode(t *testing.T) {
	name := "base"
	mock := screen_test.MockByName(name)

	node := New(mock).ToNode()
	screen_test.Helper_ToNode(t, node)

	assert.Equal(t, node.Name, name)
}

func TestHistory_Propagate(t *testing.T) {
	name := "base"
	mock := screen_test.MockByName(name)

	node := New(mock).ToNode()
	screen_test.Helper_Propagate(t, name, 0, node)
}

func TestHistory_Navigation(t *testing.T) {
	uiState := &state.UIState{}

	definition := screen.NewDefinition(
		key.ResolveDescriptors,
		key.ActionEnter,
	)

	mockBase := screen_test.MockNode{
		Name:       "base",
		Definition: &definition,
	}

	mockNext := screen_test.MockNode{
		Name:       "next",
		Definition: &definition,
		Tick: func(s *state.UIState, e screen.Event) screen.Result {
			base := mockBase.ToNode()
			return screen.ResultFromNode(&base)
		},
	}

	eventBase := screen.NewEvent(*key.NewKeyCode(key.ActionEnter))
	eventPrev := screen.NewEvent(*key.NewKeyCode(key.CustomActionPrev))
	eventNext := screen.NewEvent(*key.NewKeyCode(key.CustomActionNext))

	node := New(mockNext.ToNode()).ToNode()
	assert.Equal(t, node.Name, "next")

	result := node.Screen.Tick(uiState, eventBase)
	assert.NotNil(t, result.Node)
	assert.Equal(t, result.Node.Name, "base")

	backResult := result.Node.Screen.Tick(uiState, eventPrev)
	assert.NotNil(t, backResult.Node.Screen)
	assert.Equal(t, backResult.Node.Name, "next")

	nextResult := result.Node.Screen.Tick(uiState, eventNext)
	assert.NotNil(t, nextResult.Node.Screen)
	assert.Equal(t, nextResult.Node.Name, "base")
}

func TestHistory_ViewFooter(t *testing.T) {
	mock1 := screen_test.MockByName("mock_1")
	mock2 := screen_test.MockByName("mock_2")
	mock3 := screen_test.MockByName("mock_3")

	history := New(mock1)

	vm := history.view(*state.NewUIState())

	footer := vm.Footer.ToUnit()
	footer.Drawable.Boot()

	lines, _ := footer.Drawable.Draw(winsize.Winsize{})

	assert.Size(t, 0, lines)

	history.trail.GoTo(mock2)
	history.trail.GoTo(mock3)
	history.trail.Back()

	vm = history.view(*state.NewUIState())

	footer = vm.Footer.ToUnit()
	footer.Drawable.Boot()

	lines, _ = footer.Drawable.Draw(winsize.Winsize{
		Rows: 3,
		Cols: 30,
	})

	back := fmt.Sprintf("%s %s", history.meta.BackTag, mock1.Name)
	next := fmt.Sprintf("%s %s", history.meta.NextTag, mock3.Name)
	want := fmt.Sprintf("%s%s%s", back, history.meta.Separator, next)

	assert.Equal(t, want, text.LineToString(&lines[0]))
}
