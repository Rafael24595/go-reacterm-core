package pager

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager/predicate"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
)

func TestPagerStrategy_Integration(t *testing.T) {
	strategy := NewStrategy()

	ctx := predicate.Context{
		Page:     1,
		HasFocus: false,
	}

	pager := state.PagerContext{
		TargetPage: 1,
	}

	shouldStop := strategy.Predicate.Handler(pager, ctx)

	assert.True(t, shouldStop)
}
