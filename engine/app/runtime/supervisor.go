package runtime

import (
	"github.com/Rafael24595/go-reacterm-core/engine/build"
	"github.com/Rafael24595/go-supervisor/supervisor/result"
)

// DefaultRestartIf determines whether a supervised process should restart based on the execution result.
// It returns true in release builds to maintain application availability, and false in debug builds to allow immediate error inspection.
func DefaultRestartIf(res result.Result) bool {
	return !build.Debug
}
