package behavior

import "github.com/Rafael24595/go-reacterm-core/engine/app/screen"

type Context[T screen.Funcs] struct {
	Target Target
	Next   T
}

func NewContext[T screen.Funcs](target Target, next T) Context[T] {
	return Context[T]{
		Target: target,
		Next:   next,
	}
}
