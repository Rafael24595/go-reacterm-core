package chunk

import (
	"strings"
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/model/offset"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

func benchmarkLine(size int) line.Line {
	return line.NewBuilder().
		PushText(strings.Repeat("A", size)).
		Line()
}

func TestLine_Empty(t *testing.T) {
	src := line.Empty()

	dst := Line(src, DefaultChunk)

	assert.Equal(t, 0, dst.Size())
	assert.Equal(t, src.Hash(), dst.Hash())
}

func TestLine_ZeroLimit_ReturnsOriginal(t *testing.T) {
	src := line.FromString("golang")

	dst := Line(src, 0)

	assert.Equal(t, src.Hash(), dst.Hash())
	assert.Equal(t, src.Size(), dst.Size())
}

func TestLine_NoSplit(t *testing.T) {
	src := line.NewBuilder().
		PushText("golang").
		Line()

	dst := Line(src, 64)

	assert.Equal(t, 1, dst.Size())
	assert.Equal(t, "golang", dst.AtOrZero(0).Text())
}

func TestLine_ExactLimit(t *testing.T) {
	text := strings.Repeat("A", 64)

	src := line.NewBuilder().
		PushText(text).
		Line()

	dst := Line(src, 64)

	assert.Equal(t, 1, dst.Size())
	assert.Equal(t, text, dst.AtOrZero(0).Text())
}

func TestLine_Split(t *testing.T) {
	src := line.FromString("ziglang")

	dst := Line(src, 3)

	assert.Equal(t, 3, dst.Size())

	assert.Equal(t, "zig", dst.AtOrZero(0).Text())
	assert.Equal(t, "lan", dst.AtOrZero(1).Text())
	assert.Equal(t, "g", dst.AtOrZero(2).Text())
}

func TestLine_MultipleFragments(t *testing.T) {
	src := line.NewBuilder().
		PushText("abc").
		PushText("123456789").
		PushText("xyz").
		Line()

	dst := Line(src, 4)

	assert.Equal(t, 5, dst.Size())

	assert.Equal(t, "abc", dst.AtOrZero(0).Text())
	assert.Equal(t, "1234", dst.AtOrZero(1).Text())
	assert.Equal(t, "5678", dst.AtOrZero(2).Text())
	assert.Equal(t, "9", dst.AtOrZero(3).Text())
	assert.Equal(t, "xyz", dst.AtOrZero(4).Text())
}

func TestLine_PreserveSpec(t *testing.T) {
	style := spec.Fill(10)

	frg := frag.NewBuilder().
		AddText(strings.Repeat("A", 100)).
		AddSpec(style).
		Frag()

	src := line.NewBuilder().
		PushFrags(frg).
		Line()

	dst := Line(src, 16)

	for f := range dst.All() {
		s := f.Spec()
		assert.Equal(t, style.Hash(), s.Hash())
	}
}

func TestLine_PreserveAtom(t *testing.T) {
	frg := frag.NewBuilder().
		AddText(strings.Repeat("A", 100)).
		AddAtom(atom.Focus).
		Frag()

	src := line.NewBuilder().
		PushFrags(frg).
		Line()

	dst := Line(src, 16)

	for f := range dst.All() {
		assert.Equal(t, atom.Focus, f.Atom())
	}
}

func TestLine_PreserveOrder(t *testing.T) {
	src := line.NewBuilder().
		SetOrder(42).
		PushText(strings.Repeat("A", 100)).
		Line()

	dst := Line(src, 8)

	assert.Equal(t, uint16(42), dst.Order())
}

func TestSplitAt_NoSplit(t *testing.T) {
	src := frag.NewBuilder().
		AddText("Hello").
		Frag()

	head, tail := splitAt(src, 64)

	assert.Equal(t, "Hello", head.Text())
	assert.Nil(t, tail)
}

func TestSplitAt_Split(t *testing.T) {
	src := frag.NewBuilder().
		AddText("abcdefghij").
		Frag()

	head, tail := splitAt(src, 4)

	assert.Equal(t, "abcd", head.Text())
	assert.NotNil(t, tail)
	assert.Equal(t, "efghij", tail.Text())
}

func TestSplitAt_Unicode(t *testing.T) {
	src := frag.NewBuilder().
		AddText("áéíóú😀🚀").
		Frag()

	head, tail := splitAt(src, offset.Offset(4))

	assert.Equal(t, "áéíó", head.Text())
	assert.NotNil(t, tail)
	assert.Equal(t, "ú😀🚀", tail.Text())
}

func TestSplitAt_Emoji(t *testing.T) {
	src := frag.NewBuilder().
		AddText("😀😀😀😀😀").
		Frag()

	head, tail := splitAt(src, 3)

	assert.Equal(t, "😀😀😀", head.Text())
	assert.NotNil(t, tail)
	assert.Equal(t, "😀😀", tail.Text())
}

func BenchmarkShort(b *testing.B) {
	src := benchmarkLine(16)

	b.ReportAllocs()

	for b.Loop() {
		_ = Line(src, DefaultChunk)
	}
}

func BenchmarkExactChunk(b *testing.B) {
	src := benchmarkLine(64)

	b.ReportAllocs()

	for b.Loop() {
		_ = Line(src, DefaultChunk)
	}
}

func Benchmark1KB(b *testing.B) {
	src := benchmarkLine(1024)

	b.ReportAllocs()

	for b.Loop() {
		_ = Line(src, DefaultChunk)
	}
}

func Benchmark10KB(b *testing.B) {
	src := benchmarkLine(10 * 1024)

	b.ReportAllocs()

	for b.Loop() {
		_ = Line(src, DefaultChunk)
	}
}

func Benchmark100KB(b *testing.B) {
	src := benchmarkLine(100 * 1024)

	b.ReportAllocs()

	for b.Loop() {
		_ = Line(src, DefaultChunk)
	}
}

func BenchmarkUnicode(b *testing.B) {
	src := line.NewBuilder().
		PushText(strings.Repeat("😀", 1024)).
		Line()

	b.ReportAllocs()

	for b.Loop() {
		_ = Line(src, offset.Offset(DefaultChunk))
	}
}

func BenchmarkManyFragments(b *testing.B) {
	builder := line.NewBuilder()

	for i := 0; i < 1000; i++ {
		builder.PushText("abcdefghij")
	}

	src := builder.Line()

	b.ReportAllocs()

	for b.Loop() {
		_ = Line(src, DefaultChunk)
	}
}
