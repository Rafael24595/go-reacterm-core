package core

import (
	"context"
	"testing"
	"time"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/partial/dummy"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"

	cleaner_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/cleaner"
	screen_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/screen"
	layout_test "github.com/Rafael24595/go-reacterm-core/test/engine/layout"
	render_test "github.com/Rafael24595/go-reacterm-core/test/engine/render"
	terminal_test "github.com/Rafael24595/go-reacterm-core/test/engine/terminal"
)

func TestEngine_Lifecycle_Exit(t *testing.T) {
	trm := terminal_test.DiscardTerminal()
	lyt := layout_test.DiscardLayout()
	rnd := render_test.DiscardRender()
	cln := cleaner_test.DiscardCleaner()

	node := dummy.ToNode()

	eng := NewEngine(trm, lyt, rnd, cln, node)

	done := eng.Run()
	assert.NotNil(t, done)

	eng.Exit()

	select {
	case <-done:
	case <-time.After(1 * time.Second):
		assert.Unreachable(t, "Timeout: The done channel should be closed after calling Exit()")
	}
}

func TestEngine_Lifecycle_ContextCancel(t *testing.T) {
	trm := terminal_test.DiscardTerminal()
	lyt := layout_test.DiscardLayout()
	rnd := render_test.DiscardRender()
	cln := cleaner_test.DiscardCleaner()

	node := dummy.ToNode()

	ctx, cancel := context.WithCancel(context.Background())
	eng := NewEngine(trm, lyt, rnd, cln, node)

	done := eng.RunWithContext(ctx)

	cancel()

	select {
	case <-done:
	case <-time.After(1 * time.Second):
		assert.Unreachable(t, "Timeout: The done channel should be closed after canceling the context")
	}
}

func TestEngine_Lifecycle_ActionExitKey(t *testing.T) {
	mockTerm := terminal_test.MockTerminal{
		Keys: make(chan key.Key, 1),
	}

	trm := mockTerm.ToTerminal(winsize.New(80, 24))
	lyt := layout_test.DiscardLayout()
	rnd := render_test.DiscardRender()
	cln := cleaner_test.DiscardCleaner()

	node := dummy.ToNode()

	eng := NewEngine(trm, lyt, rnd, cln, node)
	done := eng.Run()

	mockTerm.Keys <- key.Key{Code: key.ActionExit}

	select {
	case <-done:
	case <-time.After(1 * time.Second):
		assert.Unreachable(t, "Timeout: The done channel should be closed after receiving the ActionExit key")
	}
}

func TestEngine_MutationsBlockedWhenRunning(t *testing.T) {
	trm := terminal_test.DiscardTerminal()
	lyt := layout_test.DiscardLayout()
	rnd := render_test.DiscardRender()
	cln := cleaner_test.DiscardCleaner()

	node := dummy.ToNode()

	eng := NewEngine(trm, lyt, rnd, cln, node)
	eng.Run()

	defer eng.Exit()

	assert.Panic(t, func() {
		eng.Context(context.Background())
	})

	assert.Panic(t, func() {
		eng.AddPass()
	})
}

func TestEngine_Event_ResizeTriggersRender(t *testing.T) {
	mockTerm := terminal_test.MockTerminal{
		Resize: make(chan winsize.Winsize, 1),
	}

	trm := mockTerm.ToTerminal(winsize.New(80, 24))

	renderedSize := make(chan winsize.Winsize, 1)

	composer := func(
		s *state.UIState,
		vm viewmodel.ViewModel,
		ws winsize.Winsize,
	) (*state.UIState, []line.Line) {
		select {
		case renderedSize <- ws:
		default:
		}
		return s, nil
	}

	lyt := layout_test.FromComposer(composer)
	rnd := render_test.DiscardRender()
	cln := cleaner_test.DiscardCleaner()

	node := dummy.ToNode()

	eng := NewEngine(trm, lyt, rnd, cln, node)
	eng.Run()
	defer eng.Exit()

	<-renderedSize

	newSize := winsize.New(120, 40)
	mockTerm.Resize <- newSize

	select {
	case sz := <-renderedSize:
		assert.Equal(t, newSize, sz)
	case <-time.After(1 * time.Second):
		assert.Unreachable(t, "Timeout: The resize event did not trigger a render")
	}
}

func TestEngine_Event_KeyTriggersTickAndCleanup(t *testing.T) {
	mockTerm := terminal_test.MockTerminal{
		Keys: make(chan key.Key, 1),
	}

	trm := mockTerm.ToTerminal(winsize.New(80, 24))
	lyt := layout_test.DiscardLayout()
	rnd := render_test.DiscardRender()

	node := dummy.ToNode()

	cleaned := make(chan struct{}, 1)

	cln := func(
		r screen.Result,
		s *state.UIState,
	) *state.UIState {
		select {
		case cleaned <- struct{}{}:
		default:
		}
		return s
	}

	eng := NewEngine(trm, lyt, rnd, cln, node)
	eng.Run()
	defer eng.Exit()

	mockTerm.Keys <- key.Key{Code: key.ActionRune, Rune: 'a'}

	select {
	case <-cleaned:
	case <-time.After(1 * time.Second):
		assert.Unreachable(t, "Timeout: The key event did not trigger the cleanup phase")
	}
}

func TestEngine_Navigation_DynamicNodeReplacement(t *testing.T) {
	navigated := make(chan struct{}, 1)

	mockTerm := terminal_test.MockTerminal{
		Keys: make(chan key.Key, 1),
	}

	trm := mockTerm.ToTerminal(winsize.New(80, 24))
	lyt := layout_test.DiscardLayout()
	rnd := render_test.DiscardRender()
	cln := cleaner_test.DiscardCleaner()

	node := screen_test.MockNode{
		Tick: func(s *state.UIState, e screen.Event) screen.Result {
			select {
			case navigated <- struct{}{}:
			default:
			}

			return screen.ResultFromNode(
				dummy.ToNode(),
			)
		},
	}.ToNode()

	eng := NewEngine(trm, lyt, rnd, cln, node)
	eng.Run()
	defer eng.Exit()

	mockTerm.Keys <- key.Key{Code: key.ActionRune, Rune: 'n'}

	select {
	case <-navigated:
	case <-time.After(1 * time.Second):
		assert.Unreachable(t, "Timeout: The key event did not trigger the navigation to a new node")
	}
}
