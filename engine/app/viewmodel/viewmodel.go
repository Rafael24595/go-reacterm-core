package viewmodel

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/spatial/stack"
)

// TODO: Use Screen and Units sources to manage Header and Footer.

// ViewModel represents the structural layout breakdown (Header, Kernel, Footer)
// and presentation strategies emitted by a screen.
type ViewModel struct {
	// Header represents the top section of the screen, typically used for titles or navigation.
	Header   *stack.VStackUnit
	// Kernel represents the main content area of the screen.
	Kernel   *stack.VStackUnit
	// Footer represents the bottom section of the screen, often used for status or controls.
	Footer   *stack.VStackUnit
	// Pager defines the strategy for content pagination and scrolling behavior.
	Pager    *pager.Strategy
	// Behavior encapsulates the interaction and event handling context for the screen.
	Behavior BehaviorContext
}

// New instantiates a fresh ViewModel with initialized vertical stack components
// and default pager strategies.
func New() *ViewModel {
	return &ViewModel{
		Header:   stack.NewVStack(),
		Kernel:   stack.NewVStack(),
		Footer:   stack.NewVStack(),
		Pager:    pager.NewStrategy(),
		Behavior: BehaviorContext{},
	}
}

// Clone creates a deep copy of the ViewModel structural units and its pager strategy configuration.
func (v *ViewModel) Clone() *ViewModel {
	vm := New()

	vm.Header.Push(v.Header.Units()...)
	vm.Kernel.Push(v.Kernel.Units()...)
	vm.Footer.Push(v.Footer.Units()...)
	
	if v.Pager != nil {
		vm.Pager = v.Pager.Clone()
	}

	vm.Behavior = v.Behavior

	return vm
}
