package store

import (
	"github.com/Rafael24595/go-reacterm-core/engine/commons"
	"github.com/Rafael24595/go-reacterm-core/engine/platform/clock"
)

type Argument struct {
	timestamp int64
	argument  commons.Argument
}

func newArgument(clk clock.Clock, arg any) Argument {
	return Argument{
		timestamp: clk(),
		argument:  commons.ArgumentFrom(arg),
	}
}
