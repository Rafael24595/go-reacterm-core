package processor

import (
	"strconv"
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/commons/dynamic"
	"github.com/Rafael24595/go-reacterm-core/engine/format"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
	"github.com/Rafael24595/go-reacterm-core/engine/render/styler"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

func mockAtomBold(s string) string {
	return "(b:" + s + ")"
}

func mockSpecJustifyCenter(s spec.Spec, c winsize.Cols, t format.Text) (string, bool) {
	args := s.Args()

	sz := dynamic.MapOr(args[spec.KeyJustifyCenterSize], c)
	sz = min(c, sz)

	fl := args[spec.KeyJustifyCenterText].StringOr(" ")

	size := strconv.FormatUint(uint64(sz), 10) + ":"
	filler := "'" + fl + "'" + ":"

	return "[jc:" + filler + size + t.Data + "]", true
}

func mockSpecFill(s spec.Spec, c winsize.Cols, t format.Text) (string, bool) {
	args := s.Args()

	if t.Data == "" {
		t = format.TextFromString(" ")
	}

	sz := dynamic.MapOr(args[spec.KeyFillSize], c)
	sz = min(c, sz)

	size := strconv.FormatUint(uint64(sz), 10) + ":"
	text := "'" + t.Data + "'" + ":"

	return "[fl:" + text + size + "]", true
}

func TestProcessor_EmptyInput(t *testing.T) {
	proc := New(*styler.NewAtom(), *styler.NewSpec())
	size := winsize.New(10, 80)

	t.Run("Render returns empty string for empty lines slice", func(t *testing.T) {
		assert.Equal(t, "", proc.Render(nil, size))
		assert.Equal(t, "", proc.Render([]line.Line{}, size))
	})

	t.Run("RawRender returns empty slice for empty lines", func(t *testing.T) {
		raw := proc.RawRender([]line.Line{}, size)
		assert.Equal(t, 0, len(raw))
	})
}

func TestProcessor_SingleLineProcessing(t *testing.T) {
	atoms := styler.NewAtom().
		Push(styler.AtomRule{
			Atom: atom.Bold,
			Fn:   mockAtomBold,
		})

	specs := styler.NewSpec()

	proc := New(*atoms, *specs)
	size := winsize.New(1, 40)

	t.Run("renders single fragment without specs or atoms", func(t *testing.T) {
		lne := line.FromString("Hello World")

		result := proc.Render([]line.Line{lne}, size)
		assert.Equal(t, "Hello World", result)
	})

	t.Run("applies atomic style to fragment groups", func(t *testing.T) {
		f1 := frag.NewBuilder("Bold").AddAtom(atom.Bold).Frag()
		f2 := frag.NewBuilder("Text").AddAtom(atom.Bold).Frag()
		f3 := frag.FromString(" Normal")

		lne := line.FromFrags(f1, f2, f3)

		result := proc.Render([]line.Line{lne}, size)
		assert.Equal(t, "(b:BoldText) Normal", result)
	})
}

func TestProcessor_MultiLineProcessing(t *testing.T) {
	proc := New(*styler.NewAtom(), *styler.NewSpec())
	size := winsize.New(2, 40)

	l1 := line.FromString("Line 1")
	l2 := line.FromString("Line 2")

	t.Run("RawRender returns independent line array", func(t *testing.T) {
		raw := proc.RawRender([]line.Line{l1, l2}, size)

		assert.Equal(t, 2, len(raw))
		assert.Equal(t, "Line 1", raw[0])
		assert.Equal(t, "Line 2", raw[1])
	})

	t.Run("Render joins lines with newline delimiter", func(t *testing.T) {
		rendered := proc.Render([]line.Line{l1, l2}, size)
		assert.Equal(t, "Line 1\nLine 2", rendered)
	})
}

func TestProcessor_SpecAndAtomComposition(t *testing.T) {
	atoms := styler.NewAtom().
		Push(styler.AtomRule{
			Atom: atom.Bold,
			Fn:   mockAtomBold,
		})

	specs := styler.NewSpec().
		Push(
			styler.SpecRule{
				Kind: spec.KindJustifyCenter,
				Fn:   mockSpecJustifyCenter,
			},
			styler.SpecRule{
				Kind: spec.KindFill,
				Fn:   mockSpecFill,
			},
		)

	proc := New(*atoms, *specs)
	size := winsize.New(5, 50)

	frg := frag.NewBuilder("Styled").AddAtom(atom.Bold).AddSpec(spec.Cover()).Frag()
	lne := line.NewBuilder().PushFrags(frg).AddSpec(spec.AlignCenter()).Line()

	result := proc.Render([]line.Line{lne}, size)

	assert.Equal(t, "[jc:' ':50:(b:[fl:'Styled':50:])]", result)
}
