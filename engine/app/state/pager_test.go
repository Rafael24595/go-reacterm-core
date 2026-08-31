package state

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestPagerContext_NavigationFlow(t *testing.T) {
	pager := &PagerContext{}

	pager.IncTarget()
	assert.Equal(t, 1, pager.TargetPage)
	assert.True(t, pager.modified)
	assert.False(t, pager.Synchronized)

	pager.ConfirmPage()
	assert.Equal(t, 1, pager.CurrentPage)
	assert.False(t, pager.modified)

	pager.DecTarget()
	pager.DecTarget()
	assert.Equal(t, 0, pager.TargetPage)
	assert.True(t, pager.modified)

	pager.ConfirmPage(3)
	assert.Equal(t, 3, pager.TargetPage)
	assert.Equal(t, 3, pager.CurrentPage)
	assert.False(t, pager.modified)
}
