package sink

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/commons/dynamic"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

func TestApplySinks_PaddingLeft(t *testing.T) {
	line := line.Empty().
		AddSpec(spec.JustifyRight(5, "-"))

	assert.Empty(t, line.Text)

	ApplySinks(line, 80)

	assert.False(t, line.Spec.Kind().HasAny(spec.KindJustifyRight))
	assert.Size(t, 1, line.Text)

	firstFrag := line.Text[0]
	assert.True(t, firstFrag.Spec.Kind().HasAny(spec.KindJustifyRight))
	assert.Equal(t, 5, dynamic.MapOr[winsize.Cols](firstFrag.Spec.Args()[spec.KeyJustifyRightSize], 0))
}
