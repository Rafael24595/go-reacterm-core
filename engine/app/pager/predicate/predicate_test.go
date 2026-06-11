package predicate

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
)

func TestPredicatePage(t *testing.T) {
	p := Page()

	pager := state.PagerContext{
		TargetPage: 2,
	}

	tests := []struct {
		name string
		ctx  Context
		want bool
	}{
		{
			name: "same page",
			ctx:  Context{Page: 2},
			want: true,
		},
		{
			name: "different page",
			ctx:  Context{Page: 1},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := p.Handler(pager, tt.ctx)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPredicateFocus(t *testing.T) {
	p := Focus()

	tests := []struct {
		name string
		ctx  Context
		want bool
	}{
		{
			name: "has focus",
			ctx:  Context{HasFocus: true},
			want: true,
		},
		{
			name: "no focus",
			ctx:  Context{HasFocus: false},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := p.Handler(state.PagerContext{}, tt.ctx)
			assert.Equal(t, tt.want, got)
		})
	}
}
