package behavior

import "github.com/Rafael24595/go-reacterm-core/engine/app/screen"

// Context wraps target metadata and the next function in an interception chain.
type Context[T screen.Funcs] struct {
	// Target contains metadata about the screen node being intercepted.
	Target Target
	// Next points to the next function in the execution chain.
	Next   T
}

// NewContext creates and returns a new Context instance for the given target and next function.
func NewContext[T screen.Funcs](target Target, next T) Context[T] {
	return Context[T]{
		Target: target,
		Next:   next,
	}
}
