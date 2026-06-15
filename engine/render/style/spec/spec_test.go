package spec

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-reacterm-core/engine/commons"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
)

func TestEraseSpec_DeleteExists(t *testing.T) {
	scp := Merge(
		JustifyRight(10, " "),
		Fill(100),
	)

	modified, removed := Erase(scp, KindJustifyRight)
	size := commons.Mapd[winsize.Cols](removed.args[KeyJustifyRightSize], 0)

	assert.Equal(t, KindFill, modified.kind)
	assert.NotInside(t, KeyJustifyRightSize, modified.args)
	assert.Inside(t, KeyFillSize, modified.args)

	assert.Equal(t, KindJustifyRight, removed.kind)
	assert.Equal(t, 10, size)
	assert.NotInside(t, KeyFillSize, removed.args)
}

func TestEraseSpec_DeleteNonExists(t *testing.T) {
	scp := Merge(
		JustifyRight(10, " "),
		Fill(100),
	)

	modified, removed := Erase(scp, KindTruncateLeft)

	assert.Equal(t, scp.kind, modified.kind)
	assert.Equal(t, len(scp.args), len(modified.args))

	assert.Equal(t, KindNone, removed.kind)
	assert.Equal(t, 0, len(removed.args))
}

func TestEraseSpec_DeleteMultiple(t *testing.T) {
	scp := Merge(
		JustifyRight(10, " "),
		Fill(100),
	)

	toRemove := KindJustifyRight | KindFill | KindExtendRight
	modified, removed := Erase(scp, toRemove)

	assert.Equal(t, KindNone, modified.kind)
	assert.Equal(t, 0, len(modified.args))

	assert.Equal(t, KindJustifyRight|KindFill, removed.kind)
	assert.Equal(t, 3, len(removed.args))
}
