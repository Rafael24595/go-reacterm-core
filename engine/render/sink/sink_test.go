package sink

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/commons/dynamic"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"

	text_test "github.com/Rafael24595/go-reacterm-core/test/engine/render/text"
)

func TestApplySinks_NoOp(t *testing.T) {
	inputLine := line.FromSpec(spec.Empty())
	result := ApplySinks(inputLine, 80)

	assert.Equal(
		t, text_test.LineToString(inputLine), text_test.LineToString(result),
	)
}

func TestApplySinks_PaddingLeft(t *testing.T) {
	line := line.FromSpec(
		spec.JustifyRight(5, "-"),
	)

	assert.Equal(t, 0, line.Size())

	result := ApplySinks(line, 80)

	assert.False(t, result.Spec().Kind().HasAny(spec.KindJustifyRight))
	assert.Equal(t, 1, result.Size())

	firstFrag := result.AtOrZero(0)
	assert.True(t, firstFrag.Spec().Kind().HasAny(spec.KindJustifyRight))
	assert.Equal(t, 5, dynamic.MapOr[winsize.Cols](firstFrag.Spec().Args()[spec.KeyJustifyRightSize], 0))
}

func TestApplySinks_PaddingRight(t *testing.T) {
	inputLine := line.FromSpec(spec.JustifyLeft(10, " "))

	result := ApplySinks(inputLine, 80)

	assert.False(t, result.Spec().Kind().HasAny(spec.KindJustifyLeft))
	assert.Equal(t, 1, result.Size())

	lastFrag := result.AtOrZero(result.Size() - 1)
	assert.True(t, lastFrag.Spec().Kind().HasAny(spec.KindJustifyLeft))
}

func TestApplySinks_PaddingCenter_OddAvailableSpace(t *testing.T) {
	inputLine := line.FromSpec(spec.JustifyCenter(5, " "))

	result := ApplySinks(inputLine, 80)

	assert.False(t, result.Spec().Kind().HasAny(spec.KindJustifyCenter))
	assert.Equal(t, 2, result.Size())

	leftFrag := result.AtOrZero(0)
	rightFrag := result.AtOrZero(1)

	assert.True(t, leftFrag.Spec().Kind().HasAny(spec.KindJustifyRight))
	assert.True(t, rightFrag.Spec().Kind().HasAny(spec.KindJustifyLeft))
}

func TestApplySinks_PaddingCenter_NoAvailableSpace(t *testing.T) {
	inputLine := line.FromSpec(spec.JustifyCenter(2, " "))

	result := ApplySinks(inputLine, 10)

	assert.False(t, result.Spec().Kind().HasAny(spec.KindJustifyCenter))
}
