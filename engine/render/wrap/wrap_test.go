package wrap

import (
	"strings"
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/helper/runes"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/chunk"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"

	text_test "github.com/Rafael24595/go-reacterm-core/test/engine/render/text"
)

func TestWrapOnce(t *testing.T) {
	tests := []struct {
		name         string
		cols         winsize.Cols
		line         line.Line
		expectedHead string
		expectedRest string
	}{
		{
			name: "line fits",
			cols: 20,
			line: line.FromFrags(
				frag.FromString("hello world"),
			),
			expectedHead: "hello world",
			expectedRest: "",
		},
		{
			name: "wrap by words",
			cols: 10,
			line: line.FromFrags(
				frag.FromString("hello world"),
			),
			expectedHead: "hello ",
			expectedRest: "world",
		},
		{
			name: "split long word",
			cols: 5,
			line: line.FromFrags(
				frag.FromString("abcdefghij"),
			),
			expectedHead: "abcde",
			expectedRest: "fghij",
		},
		{
			name: "split fragmented long word",
			cols: 5,
			line: line.FromFrags(
				frag.FromString("abc"),
				frag.FromString("def"),
				frag.FromString("ghi"),
			),
			expectedHead: "abcde",
			expectedRest: "fghi",
		},
		{
			name: "do not split normal word if line already has content",
			cols: 8,
			line: line.FromFrags(
				frag.FromString("hello world"),
			),
			expectedHead: "hello ",
			expectedRest: "world",
		},
		{
			name: "multiple words",
			cols: 11,
			line: line.FromFrags(
				frag.FromString("hello world foo"),
			),
			expectedHead: "hello world",
			expectedRest: " foo",
		},
		{
			name: "caret split should not affect wrapping",
			cols: 20,
			line: line.FromFrags(
				frag.FromString("supercalifra"),
				frag.FromString("gilisticexp"),
				frag.FromString("ialidocious"),
			),
			expectedHead: "supercalifragilistic",
			expectedRest: "expialidocious",
		},
		{
			name: "split long word preserves trailing words",
			cols: 5,
			line: line.FromFrags(
				frag.FromString("golang"),
				frag.FromString(" "),
				frag.FromString("zig"),
				frag.FromString(" "),
				frag.FromString("rust"),
			),
			expectedHead: "golan",
			expectedRest: "g zig rust",
		},
		{
			name: "word triggers break preserves all trailing words",
			cols: 6,
			line: line.FromFrags(
				frag.FromString("rust"),
				frag.FromString(" "),
				frag.FromString("java"),
				frag.FromString(" "),
				frag.FromString("golang"),
			),
			expectedHead: "rust ",
			expectedRest: "java golang",
		},
		{
			name: "split long word that fits exactly in next lines",
			cols: 3,
			line: line.FromFrags(
				frag.FromString("ziglang"),
				frag.FromString(" "),
				frag.FromString("rust"),
			),
			expectedHead: "zig",
			expectedRest: "lang rust",
		},
		{
			name: "Without AtmBreak: moves whole word to next line if it doesn't fit",
			cols: 6,
			line: line.FromFrags(
				frag.FromString("zig "),
				frag.FromString("golang"),
			),
			expectedHead: "zig ",
			expectedRest: "golang",
		},
		{
			name: "With AtmBreak: splits word inline to fill remaining space",
			cols: 6,
			line: line.FromFrags(
				frag.FromString("zig "),
				frag.NewBuilder().
					AddText("golang").
					AddAtom(atom.Break).
					Frag(),
			),
			expectedHead: "zig go",
			expectedRest: "lang",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			words, frags := splitLineWords(tt.line)
			layout := NewLayoutLine(tt.line, words, frags)

			head, rest := wrapOnce(tt.cols, layout)

			assert.NotNil(t, head)

			headText := text_test.LineToString(head.Line())
			assert.Equal(t, tt.expectedHead, headText)

			if tt.expectedRest != "" {
				assert.NotNil(t, rest)
				assert.Equal(t, tt.expectedRest, wordsToString(rest.words, rest.frags))
			}
		})
	}
}

func TestNormalizeLines_Integrity(t *testing.T) {
	line := line.FromString("golang ziglang 10.50 rust")

	assert.Equal(t, 1, line.Size())

	layouts := NormalizeLines(line)

	assert.Size(t, 1, layouts)
	assert.Size(t, 7, layouts[0].words)
}

func TestMaterializeEmpty(t *testing.T) {
	size := winsize.Winsize{
		Cols: 10,
	}

	placeholder := " "

	tests := []struct {
		name          string
		input         []LayoutLine
		expectedCount uint
		expectedText  string
		expectedAtom  atom.Atom
	}{
		{
			name: "ShouldMaterializeTotallyEmptyLine",
			input: []LayoutLine{
				*sourceLayout(line.Empty()),
			},
			expectedCount: 1,
			expectedText:  " ",
			expectedAtom:  atom.None,
		},
		{
			name: "ShouldNotMaterializeLineWithContent",
			input: []LayoutLine{
				*sourceLayout(
					line.FromFrags(frag.FromString("Content")),
				).pushFrags(
					frag.FromString("Content"),
				),
			},
			expectedCount: 1,
			expectedText:  "Content",
			expectedAtom:  atom.None,
		},
		{
			name: "ShouldMaterializeLineWithOnlyZeroWidthFrags",
			input: []LayoutLine{
				*sourceLayout(
					line.FromString(""),
				).pushFrags(
					frag.FromString(""),
				),
			},
			expectedCount: 2,
			expectedText:  " ",
			expectedAtom:  atom.None,
		},
		{
			name: "ShouldInheritStyleFromLastZeroWidthFrag",
			input: []LayoutLine{
				*sourceLayout(
					line.FromFrags(
						frag.FromAtom(atom.Bold),
					),
				).pushFrags(
					frag.FromAtom(atom.Bold),
				),
			},
			expectedCount: 2,
			expectedText:  " ",
			expectedAtom:  atom.Bold,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MaterializeEmpty(size, placeholder, tt.input...)

			assert.Equal(t, tt.expectedCount, got[0].Source.Size())
			assert.GreaterThan(t, 0, got[0].words)
			assert.Equal(t, tt.expectedText, text_test.LineToString(got[0].Source))

			layout := got[len(got)-1]
			word := layout.words[len(layout.words)-1]
			frag := layout.frags[word.end-1]

			assert.Equal(t, tt.expectedAtom, frag.Base.Atom())
		})
	}
}

func TestWrapLine_Simple(t *testing.T) {
	line := line.TextSpec(
		"HELLO WORLD",
		spec.AlignRight(),
	)

	lines := Line(5, line)

	expected := []string{"HELLO", " ", "WORLD"}

	assert.Size(t, len(expected), lines)

	for i, l := range lines {
		var text strings.Builder
		for f := range l.All() {
			text.WriteString(f.Text())
		}

		assert.Equal(t, expected[i], text.String())
	}
}

func TestWrapLine_Styles(t *testing.T) {
	lne := line.NewBuilder().
		PushFrags(
			frag.TextAtom("HELLO", atom.Bold),
			frag.FromString(" "),
			frag.FromString("WORLD"),
		).
		AddSpec(spec.AlignRight()).
		Line()

	lines := Line(7, lne)

	assert.Size(t, 2, lines)

	assert.Equal(t, "HELLO", lines[0].AtOrZero(0).Text())
	assert.True(t, lines[0].AtOrZero(0).Atom().HasAny(atom.Bold))

	assert.Equal(t, " ", lines[0].AtOrZero(1).Text())

	assert.Equal(t, "WORLD", lines[1].AtOrZero(0).Text())
}

func TestWrapLine_LongWord(t *testing.T) {
	txt := "HELLO WORLD FROM GOLANG"

	line := line.TextSpec(txt, spec.AlignRight())

	maxWidth := winsize.Cols(10)
	lines := Line(maxWidth, line)

	for i, l := range lines {
		text := ""
		for f := range l.All() {
			text += f.Text()
		}
		if runes.Measure(text) > maxWidth {
			t.Errorf("line %d too long: %s", i, text)
		}
	}

	totalRunes := winsize.Cols(0)
	for _, l := range lines {
		for f := range l.All() {
			totalRunes += runes.Measure(f.Text())
		}
	}
	if totalRunes != runes.Measure(txt) {
		t.Errorf("total runes mismatch")
	}
}

func TestWrapLine_MultipleFrags(t *testing.T) {
	line := line.NewBuilder().
		PushFrags(
			frag.TextAtom("HELLO", atom.Bold),
			frag.TextAtom("WORLD", atom.Bold),
			frag.FromString("GO"),
		).
		AddSpec(spec.AlignRight()).
		Line()

	maxWidth := winsize.Cols(8)
	lines := Line(maxWidth, line)

	for _, l := range lines {
		width := winsize.Cols(0)
		for f := range l.All() {
			width += runes.Measure(f.Text())
		}
		if width > maxWidth {
			t.Errorf("line exceeds maxWidth: %v", l)
		}
	}
}

func TestNextLine_Fit(t *testing.T) {
	line := line.FromString("golang")

	got, remain := NextLine(10, NormalizeLines(line))

	assert.Equal(t, "golang", text_test.LineToString(*got))
	assert.Empty(t, remain)
}

func TesNextLine_Split(t *testing.T) {
	line := line.FromString("golang")

	got, remain := NextLine(2, NormalizeLines(line))

	assert.Equal(t, "go", text_test.LineToString(*got))

	assert.Size(t, 1, remain)
	assert.Equal(t, "lang", wordsToString(remain[0].words, remain[0].frags))
}

func TesNextLine_MultiFrag(t *testing.T) {
	line := line.FromFrags(
		frag.FromString("go"),
		frag.FromString(" "),
		frag.FromString("zig"),
		frag.FromString(" "),
		frag.FromString("c++"),
	)

	got, remain := NextLine(6, NormalizeLines(line))

	assert.Equal(t, "go zig", text_test.LineToString(*got))
	assert.Size(t, 1, remain)

	assert.Equal(t, " c++", wordsToString(remain[0].words, remain[0].frags))
}

func TesNextLine_BreakLongWordSingleFrag(t *testing.T) {
	line := line.FromString("golangziglangrustlang")

	got, remain := NextLine(6, NormalizeLines(line))
	assert.Equal(t, "golang", text_test.LineToString(*got))

	assert.Equal(t, "ziglangrustlang", wordsToString(remain[0].words, remain[0].frags))
}

func TesNextLine_BreakLongWordMultipleFrags(t *testing.T) {
	line := line.FromFrags(
		frag.FromString("golang"),
		frag.FromString(" "),
		frag.FromString("zigrust"),
	)

	got, remain := NextLine(10, NormalizeLines(line))
	assert.Equal(t, "golang ", text_test.LineToString(*got))

	assert.Equal(t, "zigrust", wordsToString(remain[0].words, remain[0].frags))
}

func TestSplitFragAt(t *testing.T) {
	frg := frag.FromString("ziglang")
	wrd := newWordFrag(&frg)

	left, right := splitFragAt(wrd, 3)

	assert.NotNil(t, left)
	assert.NotNil(t, right)

	assert.NotSame(t, wrd, left)
	assert.NotSame(t, wrd, right)

	assert.NotSame(t, wrd.Base, left.Base)
	assert.NotSame(t, wrd.Base, right.Base)

	assert.Equal(t, "zig", left.Base.Text())
	assert.Equal(t, "lang", right.Base.Text())
}

func TestSplitFragAt_StartOfFrag(t *testing.T) {
	frg := frag.FromString("abcdef")
	wrd := newWordFrag(&frg)

	left, right := splitFragAt(wrd, 0)

	assert.NotNil(t, left)
	assert.NotNil(t, right)

	assert.NotSame(t, wrd, left)
	assert.NotSame(t, wrd, right)

	assert.NotSame(t, wrd.Base, left.Base)
	assert.NotSame(t, wrd.Base, right.Base)

	assert.Equal(t, "abcdef", right.Base.Text())
}

func TestSplitFragAt_EndOfFrag(t *testing.T) {
	frg := frag.FromString("abcdef")
	wrd := newWordFrag(&frg)

	left, right := splitFragAt(wrd, 6)

	assert.NotNil(t, left)
	assert.Nil(t, right)

	assert.NotSame(t, wrd, left)

	assert.NotSame(t, wrd.Base, left.Base)

	assert.Equal(t, "abcdef", left.Base.Text())
}

func BenchmarkWrapNormalize(b *testing.B) {
	wrapper := NewWrapper(
		WithProcessors(
			LineFeedProcessor,
		),
		WithSplitter(
			SplitLine,
		),
	)

	l := benchmarkLine(2000)

	b.ReportAllocs()

	for b.Loop() {
		wrapper.normalizeLines(false, l)
	}
}

func BenchmarkWrapNormalizeCached(b *testing.B) {
	wrapper := NewWrapper(
		WithProcessors(
			LineFeedProcessor,
			ChunkProcessor(chunk.DefaultChunk),
		),
		WithSplitter(
			CacheLineWords(NewFragCache()),
		),
	)

	l := benchmarkLine(2000)

	b.ReportAllocs()

	for b.Loop() {
		wrapper.normalizeLines(false, l)
	}
}

func BenchmarkWrapLine_Short(b *testing.B) {
	line := benchmarkLine(20)

	b.ReportAllocs()

	for b.Loop() {
		_ = wrapLine(80, line, nil)
	}
}

func BenchmarkWrapLine_Medium(b *testing.B) {
	line := benchmarkLine(100)

	b.ReportAllocs()

	for b.Loop() {
		_ = wrapLine(80, line, nil)
	}
}

func BenchmarkWrapLine_Long(b *testing.B) {
	line := benchmarkLine(500)

	b.ReportAllocs()

	for b.Loop() {
		_ = wrapLine(80, line, nil)
	}
}

func BenchmarkWrapLine_VeryLong(b *testing.B) {
	line := benchmarkLine(2000)

	b.ReportAllocs()

	for b.Loop() {
		_ = wrapLine(winsize.Cols(80), line, nil)
	}
}

func BenchmarkWrapOnce(b *testing.B) {
	line := line.FromFrags(
		frag.FromStrings(
			"Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
		)...,
	)

	words, frags := splitLineWords(line)

	layout := NewLayoutLine(
		line, words, frags,
	)

	b.ReportAllocs()

	for b.Loop() {
		_, _ = wrapOnce(40, layout)
	}
}

func BenchmarkWrapOnce_VeryLong(b *testing.B) {
	line := benchmarkLine(2000)

	words, frags := splitLineWords(line)

	layout := NewLayoutLine(
		line, words, frags,
	)

	b.ReportAllocs()

	for b.Loop() {
		wrapOnce(20, layout)
	}
}

func BenchmarkWrap(b *testing.B) {
	l := benchmarkLine(2000)

	b.ReportAllocs()

	for b.Loop() {
		Lines(80, l)
	}
}

func BenchmarkWrap_Resize(b *testing.B) {
	lne := benchmarkLine(20_000)

	widths := []winsize.Cols{
		120,
		100,
		80,
		60,
		40,
	}

	b.ReportAllocs()

	i := 0

	for b.Loop() {
		Lines(widths[i], lne)

		i++
		if i == len(widths) {
			i = 0
		}
	}
}
