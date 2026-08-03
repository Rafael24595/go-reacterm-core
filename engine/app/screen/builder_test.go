package screen

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"
)

func TestBuilder_BasicScreen(t *testing.T) {
	name := "home"

	node := NewBuilder().
		Name(name).
		Boot(func(u state.UIState) {}).
		Tick(func(*state.UIState, Event) Result {
			return Result{}
		}).
		View(func(state.UIState) viewmodel.ViewModel {
			return viewmodel.ViewModel{}
		}).
		ToNode()

	assert.Equal(t, name, node.Name)
	assert.Empty(t, node.Stack)
	assert.Nil(t, node.Screen.Keys)
	assert.NotNil(t, node.Screen.Boot)
	assert.NotNil(t, node.Screen.Tick)
	assert.NotNil(t, node.Screen.View)
}

func TestBuilder_WithoutKeys(t *testing.T) {
	node := NewBuilder().
		Name("home").
		WithoutKeys().
		ToNode()

	assert.NotNil(t, node.Screen.Keys)
	assert.Equal(t, 0, node.Screen.Keys().RequireKeys.Size())
}

func TestBuilder_WithNodeMeta(t *testing.T) {
	original := Node{
		Name:  "home",
		Stack: set.From("root", "main"),
	}

	node := NewBuilder().
		WithNodeMeta(original).
		ToNode()

	assert.Equal(t, original.Name, node.Name)
	assert.Inside(t, "root", node.Stack)
	assert.Inside(t, "main", node.Stack)
}

func TestBuilder_WithNodeScreen(t *testing.T) {
	original := NewBuilder().
		Boot(func(u state.UIState) {}).
		Tick(func(*state.UIState, Event) Result { return Result{} }).
		View(func(state.UIState) viewmodel.ViewModel { return viewmodel.ViewModel{} }).
		ToNode()

	node := NewBuilder().
		WithNodeScreen(original).
		ToNode()

	assert.NotNil(t, node.Screen.Boot)
	assert.NotNil(t, node.Screen.Tick)
	assert.NotNil(t, node.Screen.View)
}

func TestBuilder_WithNode(t *testing.T) {
	name := "home"
	original := NewBuilder().
		Name(name).
		NameToStack().
		Boot(func(u state.UIState) {}).
		Tick(func(*state.UIState, Event) Result { return Result{} }).
		View(func(state.UIState) viewmodel.ViewModel { return viewmodel.ViewModel{} }).
		ToNode()

	node := NewBuilder().
		WithNode(original).
		ToNode()

	assert.Equal(t, name, node.Name)
	assert.Inside(t, name, node.Stack)

	assert.NotNil(t, node.Screen.Boot)
	assert.NotNil(t, node.Screen.Tick)
	assert.NotNil(t, node.Screen.View)
}

func TestBuilder_NameAsStack(t *testing.T) {
	name := "home"

	node := NewBuilder().
		NameAsStack(name).
		ToNode()

	assert.Equal(t, name, node.Name)
	assert.Inside(t, name, node.Stack)
}

func TestBuilder_NameToStack(t *testing.T) {
	name := "home"

	node := NewBuilder().
		Name(name).
		NameToStack().
		ToNode()

	assert.Inside(t, name, node.Stack)
}

func TestBuilder_IncompleteScreen(t *testing.T) {
	node := NewBuilder().Name("home").ToNode()

	assert.Nil(t, node.Screen.Boot)
	assert.Nil(t, node.Screen.Keys)
	assert.Nil(t, node.Screen.Tick)
	assert.Nil(t, node.Screen.View)
}
