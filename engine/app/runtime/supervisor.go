package runtime

import (
	"github.com/Rafael24595/go-reacterm-core/engine/build"
	"github.com/Rafael24595/go-supervisor/supervisor/result"
)

func DefaultRestartIf(res result.Result) bool {
	return !build.Debug
}
