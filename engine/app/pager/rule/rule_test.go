package rule

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
)

func TestKinds(t *testing.T) {
	assert.Equal(t, KindPage, OnPage().Kind)
	assert.Equal(t, KindFocus, OnFocus().Kind)
}

func TestOnPage(t *testing.T) {
	rule := OnPage()

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
			got := rule.Handler(pager, tt.ctx)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestOnStart(t *testing.T) {
	r := OnStart()
	assert.True(t, r.Handler(state.PagerContext{}, Context{}))
}

func TestOnEnd(t *testing.T) {
	r := OnEnd()
	assert.False(t, r.Handler(state.PagerContext{}, Context{}))
}

func TestOnFocus(t *testing.T) {
	rule := OnFocus()

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
			got := rule.Handler(state.PagerContext{}, tt.ctx)
			assert.Equal(t, tt.want, got)
		})
	}
}
