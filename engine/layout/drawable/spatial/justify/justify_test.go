package justify

import (
	"strings"
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/commons/dynamic"
	"github.com/Rafael24595/go-reacterm-core/engine/format"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"

	drawable_test "github.com/Rafael24595/go-reacterm-core/test/engine/layout/drawable"
	text_test "github.com/Rafael24595/go-reacterm-core/test/engine/render/text"
)

func renderFrags(frags []frag.Frag) string {
	var s strings.Builder
	for _, f := range frags {
		s.WriteString(f.Text())

		count := winsize.Cols(0)
		ok := false
		if args, ex := f.Spec().Args()[spec.KeyExtendRightSize]; ex {
			ok = true
			count = dynamic.MapOr[winsize.Cols](args, 0)
		}

		if args, ex := f.Spec().Args()[spec.KeyJustifyLeftSize]; ex {
			ok = true
			count = dynamic.MapOr[winsize.Cols](args, 0)
		}

		if ok {
			for range count {
				s.WriteString(" ")
			}
		}
	}

	return s.String()
}

func renderLine(cols winsize.Cols, mode style.Justify, line line.Line) string {
	filler := marker.DefaultPaddingText

	text := format.TextFromString(
		renderFrags(line.Text),
	)

	switch mode {
	case style.JustifyStart:
		return format.JustifyLeft(cols, text, filler)
	case style.JustifyEnd:
		return format.JustifyRight(cols, text, filler)
	case style.JustifyCenter, style.JustifyAround, style.JustifyEvenly:
		return format.JustifyCenter(cols, text, filler)
	}

	return text.Data
}

func TestJustify_UnitBasicSuite(t *testing.T) {
	unit := UnitFromFrags([]frag.Frag{})
	drawable_test.Test_UnitBasicSuite(t, unit)
}

func TestAddGaps_SingleFrag(t *testing.T) {
	frags := frag.FromStrings(
		"abc",
	)

	for _, mode := range []style.Justify{
		style.JustifyStart, style.JustifyEnd, style.JustifyCenter,
		style.JustifyBetween, style.JustifyAround, style.JustifyEvenly,
	} {
		result := addGaps(10, frags, 3, mode)

		assert.Size(t, 1, result)
		assert.Equal(t, "abc", result[0].Text())
		assert.Equal(t, spec.KindNone, result[0].Spec().Kind())
	}
}

func TestAddGaps_IntercalatedSpaces(t *testing.T) {
	frags := frag.FromStrings(
		"aa", "bb", "cc",
	)

	for _, mode := range []style.Justify{style.JustifyStart, style.JustifyEnd, style.JustifyCenter} {
		result := addGaps(10, frags, 6, mode)

		assert.Size(t, 5, result)
		assert.Equal(t, "aa bb cc", text_test.FragsToString(result))

		assert.Equal(t, spec.KindNone, result[0].Spec().Kind())
		assert.Equal(t, spec.KindNone, result[2].Spec().Kind())
		assert.Equal(t, spec.KindNone, result[4].Spec().Kind())
	}
}

func TestAddGaps_Between(t *testing.T) {
	frags := frag.FromStrings(
		"aa", "bb", "cc",
	)

	result := addGaps(10, frags, 6, style.JustifyBetween)

	assert.Size(t, 5, result)

	assert.Equal(t, 2, dynamic.MapOr[winsize.Cols](result[1].Spec().Args()[spec.KeyJustifyLeftSize], 0))
	assert.Equal(t, 2, dynamic.MapOr[winsize.Cols](result[3].Spec().Args()[spec.KeyJustifyLeftSize], 0))
	assert.Equal(t, spec.KindNone, result[2].Spec().Kind())

	assert.Equal(t, "aa  bb  cc", renderFrags(result))
}

func TestAddGaps_Around(t *testing.T) {
	frags := frag.FromStrings(
		"aa", "bb", "cc",
	)

	result := addGaps(11, frags, 6, style.JustifyAround)

	assert.Size(t, 5, result)

	assert.Equal(t, 2, dynamic.MapOr[winsize.Cols](result[1].Spec().Args()[spec.KeyJustifyLeftSize], 0))
	assert.Equal(t, 1, dynamic.MapOr[winsize.Cols](result[3].Spec().Args()[spec.KeyJustifyLeftSize], 0))
	assert.Equal(t, spec.KindNone, result[2].Spec().Kind())

	assert.Equal(t, "aa  bb cc", renderFrags(result))
}

func TestAddGaps_Overflow_Start(t *testing.T) {
	frags := frag.FromStrings(
		"aaaa", "bbbb",
	)

	result := addGaps(5, frags, 8, style.JustifyStart)
	assert.Equal(t, "aaaa bbbb", text_test.FragsToString(result))
}

func TestAddGaps_DoesNotMutateOriginal(t *testing.T) {
	frags := frag.FromStrings(
		"aa", "bb",
	)

	_ = addGaps(10, frags, 4, style.JustifyBetween)

	for _, f := range frags {
		assert.Equal(t, spec.KindNone, f.Spec().Kind())
	}
}

func TestJustifyLine_Start(t *testing.T) {
	frags := frag.FromStrings("aa", "bb", "cc")
	line := justifyLine(10, frags, 6, style.JustifyStart)

	assert.Size(t, 5, line.Text)
	assert.True(t, line.Spec.Kind().HasAny(spec.KindJustifyLeft))

	assert.Equal(t, "aa bb cc  ", renderLine(10, style.JustifyStart, line))
}

func TestJustifyLine_End(t *testing.T) {
	frags := frag.FromStrings("aa", "bb", "cc")
	line := justifyLine(10, frags, 6, style.JustifyEnd)

	assert.Size(t, 5, line.Text)
	assert.True(t, line.Spec.Kind().HasAny(spec.KindJustifyRight))

	assert.Equal(t, "  aa bb cc", renderLine(10, style.JustifyEnd, line))
}

func TestJustifyLine_Center(t *testing.T) {
	frags := frag.FromStrings("aa", "bb", "cc")
	line := justifyLine(10, frags, 6, style.JustifyCenter)

	assert.Size(t, 5, line.Text)
	assert.True(t, line.Spec.Kind().HasAny(spec.KindJustifyCenter))

	assert.Equal(t, " aa bb cc ", renderLine(10, style.JustifyCenter, line))
}

func TestJustifyLine_Between(t *testing.T) {
	frags := frag.FromStrings("aa", "bb", "cc")
	line := justifyLine(10, frags, 6, style.JustifyBetween)

	assert.Size(t, 5, line.Text)
	assert.True(t, line.Spec.Kind().HasNone(spec.KindJustifyRight|spec.KindJustifyLeft|spec.KindJustifyCenter))

	assert.Equal(t, "aa  bb  cc", renderLine(10, style.JustifyBetween, line))
}

func TestJustifyLine_Around(t *testing.T) {
	frags := frag.FromStrings("aa", "bb", "cc")
	line := justifyLine(18, frags, 6, style.JustifyAround)

	assert.Size(t, 5, line.Text)
	assert.True(t, line.Spec.Kind().HasAny(spec.KindJustifyCenter))

	assert.Equal(t, "   aa   bb   cc   ", renderLine(18, style.JustifyAround, line))
}

func TestJustifyLine_Evenly(t *testing.T) {
	frags := frag.FromStrings("aa", "bb", "cc")
	line := justifyLine(18, frags, 6, style.JustifyEvenly)

	assert.Size(t, 5, line.Text)
	assert.True(t, line.Spec.Kind().HasAny(spec.KindJustifyCenter))

	assert.Equal(t, "  aa    bb    cc  ", renderLine(18, style.JustifyEvenly, line))
}
