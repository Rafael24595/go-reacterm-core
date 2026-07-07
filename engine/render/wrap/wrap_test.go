package wrap

import (
	"strings"
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/helper/runes"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
	
	text_test "github.com/Rafael24595/go-reacterm-core/test/engine/render/text"
)

func assembleLines(t *testing.T, lines ...text.Line) string {
	t.Helper()

	var sb strings.Builder

	for i, l := range lines {
		_, err := sb.WriteString(
			text_test.LineToString(&l),
		)

		assert.Nil(t, err)

		if i < len(lines)-1 {
			_, err := sb.WriteString("\n")
			assert.Nil(t, err)
		}
	}

	return sb.String()
}

func benchmarkLine(words int) text.Line {
	builder := strings.Builder{}

	for range words {
		builder.WriteString("Lorem ")
	}

	line := text.NewLine("")
	line.PushFrags(
		*text.NewFrag(builder.String()),
	)

	return *line
}

func TestWrapOnce(t *testing.T) {
	tests := []struct {
		name         string
		cols         winsize.Cols
		line         *text.Line
		expectedHead string
		expectedRest string
	}{
		{
			name: "line fits",
			cols: 20,
			line: text.LineFromFrags(
				*text.NewFrag("hello world"),
			),
			expectedHead: "hello world",
			expectedRest: "",
		},
		{
			name: "wrap by words",
			cols: 10,
			line: text.LineFromFrags(
				*text.NewFrag("hello world"),
			),
			expectedHead: "hello ",
			expectedRest: "world",
		},
		{
			name: "split long word",
			cols: 5,
			line: text.LineFromFrags(
				*text.NewFrag("abcdefghij"),
			),
			expectedHead: "abcde",
			expectedRest: "fghij",
		},
		{
			name: "split fragmented long word",
			cols: 5,
			line: text.LineFromFrags(
				*text.NewFrag("abc"),
				*text.NewFrag("def"),
				*text.NewFrag("ghi"),
			),
			expectedHead: "abcde",
			expectedRest: "fghi",
		},
		{
			name: "do not split normal word if line already has content",
			cols: 8,
			line: text.LineFromFrags(
				*text.NewFrag("hello world"),
			),
			expectedHead: "hello ",
			expectedRest: "world",
		},
		{
			name: "multiple words",
			cols: 11,
			line: text.LineFromFrags(
				*text.NewFrag("hello world foo"),
			),
			expectedHead: "hello world",
			expectedRest: " foo",
		},
		{
			name: "caret split should not affect wrapping",
			cols: 20,
			line: text.LineFromFrags(
				*text.NewFrag("supercalifra"),
				*text.NewFrag("gilisticexp"),
				*text.NewFrag("ialidocious"),
			),
			expectedHead: "supercalifragilistic",
			expectedRest: "expialidocious",
		},
		{
			name: "split long word preserves trailing words",
			cols: 5,
			line: text.LineFromFrags(
				*text.NewFrag("golang"),
				*text.NewFrag(" "),
				*text.NewFrag("zig"),
				*text.NewFrag(" "),
				*text.NewFrag("rust"),
			),
			expectedHead: "golan",
			expectedRest: "g zig rust",
		},
		{
			name: "word triggers break preserves all trailing words",
			cols: 6,
			line: text.LineFromFrags(
				*text.NewFrag("rust"),
				*text.NewFrag(" "),
				*text.NewFrag("java"),
				*text.NewFrag(" "),
				*text.NewFrag("golang"),
			),
			expectedHead: "rust ",
			expectedRest: "java golang",
		},
		{
			name: "split long word that fits exactly in next lines",
			cols: 3,
			line: text.LineFromFrags(
				*text.NewFrag("ziglang"),
				*text.NewFrag(" "),
				*text.NewFrag("rust"),
			),
			expectedHead: "zig",
			expectedRest: "lang rust",
		},
		{
			name: "Without AtmBreak: moves whole word to next line if it doesn't fit",
			cols: 6,
			line: text.LineFromFrags(
				*text.NewFrag("zig "),
				*text.NewFrag("golang"),
			),
			expectedHead: "zig ",
			expectedRest: "golang",
		},
		{
			name: "With AtmBreak: splits word inline to fill remaining space",
			cols: 6,
			line: text.LineFromFrags(
				*text.NewFrag("zig "),
				*text.NewFrag("golang").
					AddAtom(atom.Break),
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

			headText := text_test.LineToString(head)
			assert.Equal(t, tt.expectedHead, headText)

			if tt.expectedRest != "" {
				assert.NotNil(t, rest)
				assert.Equal(t, tt.expectedRest, wordsToString(rest.words, rest.frags))
			}
		})
	}
}

func TestNormalizeLines_Integrity(t *testing.T) {
	line := text.NewLine("golang ziglang 10.50 rust")

	assert.Size(t, 1, line.Text)

	layouts := NormalizeLines(*line)

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
		expectedCount int
		expectedText  string
		expectedAtom  atom.Atom
	}{
		{
			name: "ShouldMaterializeTotallyEmptyLine",
			input: []LayoutLine{
				*sourceLayout(text.EmptyLine()),
			},
			expectedCount: 1,
			expectedText:  " ",
			expectedAtom:  atom.None,
		},
		{
			name: "ShouldNotMaterializeLineWithContent",
			input: []LayoutLine{
				*sourceLayout(
					text.LineFromFrags(*text.NewFrag("Content")),
				).pushFrags(
					*text.NewFrag("Content"),
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
					text.NewLine(""),
				).pushFrags(
					*text.NewFrag(""),
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
					text.LineFromFrags(
						*text.NewFrag("").AddAtom(atom.Bold),
					),
				).pushFrags(
					*text.NewFrag("").AddAtom(atom.Bold),
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

			assert.Size(t, tt.expectedCount, got[0].Source.Text)
			assert.GreaterThan(t, 0, len(got[0].words))
			assert.Equal(t, tt.expectedText, text_test.LineToString(got[0].Source))

			layout := got[len(got)-1]
			word := layout.words[len(layout.words)-1]
			frag := layout.frags[word.end-1]

			assert.Equal(t, tt.expectedAtom, frag.Base.Atom)
		})
	}
}

func TestWrapLine_Simple(t *testing.T) {
	line := text.NewLine(
		"HELLO WORLD",
		spec.AlignRight(),
	)

	lines := Line(5, line)

	expected := []string{"HELLO", " ", "WORLD"}

	assert.Size(t, len(expected), lines)

	for i, l := range lines {
		var text strings.Builder
		for _, f := range l.Text {
			text.WriteString(f.Text)
		}

		assert.Equal(t, expected[i], text.String())
	}
}

func TestWrapLine_Styles(t *testing.T) {
	line := text.LineFromFrags(
		*text.NewFrag("HELLO").AddAtom(atom.Bold),
		*text.NewFrag(" "),
		*text.NewFrag("WORLD"),
	).SetSpec(spec.AlignRight())

	lines := Line(7, line)

	assert.Equal(t, 2, len(lines))

	assert.Equal(t, "HELLO", lines[0].Text[0].Text)
	assert.True(t, lines[0].Text[0].Atom.HasAny(atom.Bold))

	assert.Equal(t, " ", lines[0].Text[1].Text)

	assert.Equal(t, "WORLD", lines[1].Text[0].Text)
}

func TestWrapLine_LongWord(t *testing.T) {
	txt := "HELLO WORLD FROM GOLANG"

	line := text.NewLine(txt,
		spec.AlignRight(),
	)

	maxWidth := winsize.Cols(10)
	lines := Line(maxWidth, line)

	for i, l := range lines {
		text := ""
		for _, f := range l.Text {
			text += f.Text
		}
		if runes.Measure(text) > maxWidth {
			t.Errorf("line %d too long: %s", i, text)
		}
	}

	totalRunes := winsize.Cols(0)
	for _, l := range lines {
		for _, f := range l.Text {
			totalRunes += runes.Measure(f.Text)
		}
	}
	if totalRunes != runes.Measure(txt) {
		t.Errorf("total runes mismatch")
	}
}

func TestWrapLine_MultipleFrags(t *testing.T) {
	line := text.LineFromFrags(
		*text.NewFrag("HELLO").AddAtom(atom.Bold),
		*text.NewFrag("WORLD").AddAtom(atom.Bold),
		*text.NewFrag("GO"),
	).SetSpec(spec.AlignRight())

	maxWidth := winsize.Cols(8)
	lines := Line(maxWidth, line)

	for _, l := range lines {
		width := winsize.Cols(0)
		for _, f := range l.Text {
			width += runes.Measure(f.Text)
		}
		if width > maxWidth {
			t.Errorf("line exceeds maxWidth: %v", l)
		}
	}
}

func TestNextLine_Fit(t *testing.T) {
	line := text.NewLine("golang")

	got, remain := NextLine(10, NormalizeLines(*line))

	assert.Equal(t, "golang", text_test.LineToString(got))
	assert.Size(t, 0, remain)
}

func TesNextLine_Split(t *testing.T) {
	line := text.NewLine("golang")

	got, remain := NextLine(2, NormalizeLines(*line))

	assert.Equal(t, "go", text_test.LineToString(got))

	assert.Size(t, 1, remain)
	assert.Equal(t, "lang", wordsToString(remain[0].words, remain[0].frags))
}

func TesNextLine_MultiFrag(t *testing.T) {
	line := text.LineFromFrags(
		*text.NewFrag("go"),
		*text.NewFrag(" "),
		*text.NewFrag("zig"),
		*text.NewFrag(" "),
		*text.NewFrag("c++"),
	)

	got, remain := NextLine(6, NormalizeLines(*line))

	assert.Equal(t, "go zig", text_test.LineToString(got))
	assert.Size(t, 1, remain)

	assert.Equal(t, " c++", wordsToString(remain[0].words, remain[0].frags))
}

func TesNextLine_BreakLongWordSingleFrag(t *testing.T) {
	line := text.NewLine("golangziglangrustlang")

	got, remain := NextLine(6, NormalizeLines(*line))
	assert.Equal(t, "golang", text_test.LineToString(got))

	assert.Equal(t, "ziglangrustlang", wordsToString(remain[0].words, remain[0].frags))
}

func TesNextLine_BreakLongWordMultipleFrags(t *testing.T) {
	line := text.LineFromFrags(
		*text.NewFrag("golang"),
		*text.NewFrag(" "),
		*text.NewFrag("zigrust"),
	)

	got, remain := NextLine(10, NormalizeLines(*line))
	assert.Equal(t, "golang ", text_test.LineToString(got))

	assert.Equal(t, "zigrust", wordsToString(remain[0].words, remain[0].frags))
}

func TestSplitLineFeeds(t *testing.T) {
	tests := []struct {
		name         string
		input        *text.Line
		expectedSize int
		expectedText string
		expecteFrags []int
	}{
		{
			name: "WithoutLineFeed",
			input: text.EmptyLine().PushFrags(
				*text.NewFrag("Hello Golang"),
			),
			expectedSize: 1,
			expectedText: "Hello Golang",
			expecteFrags: []int{1},
		},
		{
			name: "SingleLineFeed",
			input: text.EmptyLine().PushFrags(
				*text.NewFrag("Golang\nZiglang"),
			),
			expectedSize: 2,
			expectedText: "Golang\nZiglang",
			expecteFrags: []int{1, 1},
		},
		{
			name: "LineFeedBetweenFrags",
			input: text.EmptyLine().PushFrags(
				*text.NewFrag("Rust"),
				*text.NewFrag("\nZig"),
			),
			expectedSize: 2,
			expectedText: "Rust\nZig",
			expecteFrags: []int{1, 1},
		},
		{
			name: "MultipleLineFeedWithEmptyLine",
			input: text.EmptyLine().PushFrags(
				*text.NewFrag("Go\n\nC++"),
			),
			expectedSize: 3,
			expectedText: "Go\n\nC++",
			expecteFrags: []int{1, 0, 1},
		},
		{
			name: "LineFeedAtEnd",
			input: text.EmptyLine().PushFrags(
				*text.NewFrag("Rust\n"),
			),
			expectedSize: 2,
			expectedText: "Rust\n",
			expecteFrags: []int{1, 0},
		},
		{
			name: "LineFeedWithCarriageReturn",
			input: text.EmptyLine().PushFrags(
				*text.NewFrag("Zig\r\nGolang"),
			),
			expectedSize: 2,
			expectedText: "Zig\nGolang",
			expecteFrags: []int{1, 1},
		},
		{
			name: "CarriageReturn",
			input: text.EmptyLine().PushFrags(
				*text.NewFrag("Java\rElixir"),
			),
			expectedSize: 2,
			expectedText: "Java\nElixir",
			expecteFrags: []int{1, 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := splitLineFeeds(tt.input, false)

			assert.Size(t, tt.expectedSize, got)
			assert.Equal(t, tt.expectedText, assembleLines(t, got...))

			for i, v := range got {
				assert.Size(t, int(tt.expecteFrags[i]), v.Text)
			}
		})
	}
}

func TestSplitLineFeeds_Ordering(t *testing.T) {
	tests := []struct {
		name           string
		input          *text.Line
		orderFlag      bool
		expectedOrders []uint16
	}{
		{
			name:           "ShouldNotSetOrderIfFlagIsFalse",
			input:          text.NewLine("Line1\nLine2"),
			orderFlag:      false,
			expectedOrders: []uint16{0, 0},
		},
		{
			name:           "ShouldStartFromOneIfOrderIsZero",
			input:          text.NewLine("Line1\nLine2\nLine3"),
			orderFlag:      true,
			expectedOrders: []uint16{1, 2, 3},
		},
		{
			name: "ShouldResumeFromExistingOrder",
			input: text.NewLine("PartA\nPartB").
				SetOrder(10),
			orderFlag:      true,
			expectedOrders: []uint16{10, 11},
		},
		{
			name: "ShouldHandleMultipleFragsWithOrder",
			input: text.LineFromFrags(
				*text.NewFrag("A"),
				*text.NewFrag("\nB\n"),
				*text.NewFrag("C"),
			),
			orderFlag:      true,
			expectedOrders: []uint16{1, 2, 3},
		},
		{
			name: "ShouldHandleMultipleFragsWithOrder",
			input: text.LineFromFrags(
				*text.NewFrag("A"),
				*text.NewFrag("\nB\n"),
				*text.NewFrag("C"),
			),
			orderFlag:      true,
			expectedOrders: []uint16{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := splitLineFeeds(tt.input, tt.orderFlag)

			assert.Equal(t, len(tt.expectedOrders), len(got), "Result size mismatch")

			for i, line := range got {
				assert.Equal(t, tt.expectedOrders[i], line.Order, "Order mismatch at index %d", i)
			}
		})
	}
}

func TestSplitFragAt_EndOfFrag(t *testing.T) {
	frag := newWordFrag(
		text.NewFrag("abcdef"),
	)

	left, right := splitFragAt(frag, 6)

	assert.NotNil(t, left)
	assert.Nil(t, right)

	assert.Equal(t, "abcdef", left.Base.Text)
}

func TestSplitFragAt_EmptyRestNeverCreated(t *testing.T) {
	frag := newWordFrag(
		text.NewFrag("abc"),
	)

	_, right := splitFragAt(frag, 3)

	assert.Nil(t, right)
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
	line := text.LineFromFrags(
		text.FragsFromString(
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

	words, frags := splitLineWords(&line)

	layout := NewLayoutLine(
		&line, words, frags,
	)

	b.ReportAllocs()

	for b.Loop() {
		wrapOnce(20, layout)
	}
}
