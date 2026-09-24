package clock

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

)

func TestUnixMilliClock(t *testing.T) {
	ts1 := UnixMilliClock()
	assert.True(t, ts1 > 0)
}

func TestGlobalCounterClock(t *testing.T) {
	resetGlobalCounter()

	c1 := GlobalCounterClock()
	c2 := GlobalCounterClock()

	assert.Equal(t, int64(1), c1)
	assert.Equal(t, int64(2), c2)
}