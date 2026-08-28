package pass

import (
	"fmt"
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"

	screen_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/screen"
)

func TestValidateStructure_ValidNode(t *testing.T) {
	node := screen_test.MockNode{
		Name: "home",
		Boot: func(state.UIState) {},
		Tick: func(*state.UIState, screen.Event) screen.Result {
			return screen.Result{}
		},
		View: func(state.UIState) viewmodel.ViewModel {
			return viewmodel.ViewModel{}
		},
	}.ToNode()

	_, err := ValidateStructure(node)

	assert.Nil(t, err)
}

func TestValidateStructure_EmptyName(t *testing.T) {
	node := screen.Node{
		Screen: screen.Screen{
			Keys: func() screen.Definition {
				return screen.Definition{}
			},
			Tick: func(*state.UIState, screen.Event) screen.Result {
				return screen.Result{}
			},
			View: func(state.UIState) viewmodel.ViewModel {
				return viewmodel.ViewModel{}
			},
		},
	}

	_, err := ValidateStructure(node)

	assert.NotNil(t, err)
	assert.Equal(t, err_name, err.Error())
}

func TestValidateStructure_NilKeys(t *testing.T) {
	name := "home"

	node := screen.Node{
		Name: name,
		Screen: screen.Screen{
			Boot: func(u state.UIState) {},
			Tick: func(*state.UIState, screen.Event) screen.Result {
				return screen.Result{}
			},
			View: func(state.UIState) viewmodel.ViewModel {
				return viewmodel.ViewModel{}
			},
		},
	}

	_, err := ValidateStructure(node)

	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf(errf_keys, name), err.Error())
}

func TestValidateStructure_NilBoot(t *testing.T) {
	name := "home"

	node := screen.Node{
		Name: name,
		Screen: screen.Screen{
			Keys: func() screen.Definition {
				return screen.Definition{}
			},
			View: func(state.UIState) viewmodel.ViewModel {
				return viewmodel.ViewModel{}
			},
		},
	}

	_, err := ValidateStructure(node)

	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf(errf_boot, name), err.Error())
}

func TestValidateStructure_NilTick(t *testing.T) {
	name := "home"

	node := screen.Node{
		Name: name,
		Screen: screen.Screen{
			Boot: func(u state.UIState) {},
			Keys: func() screen.Definition {
				return screen.Definition{}
			},
			View: func(state.UIState) viewmodel.ViewModel {
				return viewmodel.ViewModel{}
			},
		},
	}

	_, err := ValidateStructure(node)

	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf(errf_tick, name), err.Error())
}

func TestValidateStructure_NilView(t *testing.T) {
	name := "home"

	node := screen.Node{
		Name: name,
		Screen: screen.Screen{
			Boot: func(u state.UIState) {},
			Keys: func() screen.Definition {
				return screen.Definition{}
			},
			Tick: func(*state.UIState, screen.Event) screen.Result {
				return screen.Result{}
			},
		},
	}

	_, err := ValidateStructure(node)

	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf(errf_view, name), err.Error())
}

func TestValidateStructure_CycleDetected(t *testing.T) {
	parent := screen.NewBuilder().
		WithClock(func() int64 {
			return 0
		}).
		WithNode(screen_test.DummyNode).
		Name("parent")

	child := screen.NewBuilder().
		WithClock(func() int64 {
			return 1
		}).
		WithNode(screen_test.DummyNode).
		Name("child")

	child.Children(parent.ToNode())
	parent.Children(child.ToNode())

	_, err := ValidateStructure(parent.ToNode())

	println(err.Error())

	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf(errf_cycle, name), err.Error())
}
