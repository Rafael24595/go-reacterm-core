package store

import (
	"github.com/Rafael24595/go-reacterm-core/engine/commons/dynamic"
	"github.com/Rafael24595/go-reacterm-core/engine/platform/clock"
)

// TODO: Improve nomenclature.
type Argument struct {
	timestamp int64
	argument  dynamic.Value
}

func newArgument(clk clock.Clock, arg any) Argument {
	return Argument{
		timestamp: clk(),
		argument:  dynamic.From(arg),
	}
}
