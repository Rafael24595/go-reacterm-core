package processor

import (
	"strings"
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

func TestLineFeed(t *testing.T) {
	tests := []struct {
		name         string
		input        line.Line
		expectedSize int
		expectedText string
		expecteFrags []uint
	}{
		{
			name: "WithoutLineFeed",
			input: line.FromFrags(
				frag.FromString("Hello Golang"),
			),
			expectedSize: 1,
			expectedText: "Hello Golang",
			expecteFrags: []uint{1},
		},
		{
			name: "SingleLineFeed",
			input: line.FromFrags(
				frag.FromString("Golang\nZiglang"),
			),
			expectedSize: 2,
			expectedText: "Golang\nZiglang",
			expecteFrags: []uint{1, 1},
		},
		{
			name: "LineFeedBetweenFrags",
			input: line.FromFrags(
				frag.FromString("Rust"),
				frag.FromString("\nZig"),
			),
			expectedSize: 2,
			expectedText: "Rust\nZig",
			expecteFrags: []uint{1, 1},
		},
		{
			name: "MultipleLineFeedWithEmptyLine",
			input: line.FromFrags(
				frag.FromString("Go\n\nC++"),
			),
			expectedSize: 3,
			expectedText: "Go\n\nC++",
			expecteFrags: []uint{1, 0, 1},
		},
		{
			name: "LineFeedAtEnd",
			input: line.FromFrags(
				frag.FromString("Rust\n"),
			),
			expectedSize: 2,
			expectedText: "Rust\n",
			expecteFrags: []uint{1, 0},
		},
		{
			name: "LineFeedWithCarriageReturn",
			input: line.FromFrags(
				frag.FromString("Zig\r\nGolang"),
			),
			expectedSize: 2,
			expectedText: "Zig\nGolang",
			expecteFrags: []uint{1, 1},
		},
		{
			name: "CarriageReturn",
			input: line.FromFrags(
				frag.FromString("Java\rElixir"),
			),
			expectedSize: 2,
			expectedText: "Java\nElixir",
			expecteFrags: []uint{1, 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LineFeed(false, tt.input)

			assert.Size(t, tt.expectedSize, got)
			assert.Equal(t, tt.expectedText, assembleLines(t, got...))

			for i, v := range got {
				assert.Equal(t, tt.expecteFrags[i], v.Size())
			}
		})
	}
}

func TestLineFeed_Ordering(t *testing.T) {
	tests := []struct {
		name           string
		input          line.Line
		orderFlag      bool
		expectedOrders []uint16
	}{
		{
			name:           "ShouldNotSetOrderIfFlagIsFalse",
			input:          line.FromString("Line1\nLine2"),
			orderFlag:      false,
			expectedOrders: []uint16{0, 0},
		},
		{
			name:           "ShouldStartFromOneIfOrderIsZero",
			input:          line.FromString("Line1\nLine2\nLine3"),
			orderFlag:      true,
			expectedOrders: []uint16{1, 2, 3},
		},
		{
			name:           "ShouldResumeFromExistingOrder",
			input:          line.TextOrder(10, "PartA\nPartB"),
			orderFlag:      true,
			expectedOrders: []uint16{10, 11},
		},
		{
			name: "ShouldHandleMultipleFragsWithOrder",
			input: line.FromFrags(
				frag.FromString("A"),
				frag.FromString("\nB\n"),
				frag.FromString("C"),
			),
			orderFlag:      true,
			expectedOrders: []uint16{1, 2, 3},
		},
		{
			name: "ShouldHandleMultipleFragsWithOrder",
			input: line.FromFrags(
				frag.FromString("A"),
				frag.FromString("\nB\n"),
				frag.FromString("C"),
			),
			orderFlag:      true,
			expectedOrders: []uint16{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LineFeed(tt.orderFlag, tt.input)

			assert.Equal(t, len(tt.expectedOrders), len(got), "Result size mismatch")

			for i, line := range got {
				assert.Equal(
					t, tt.expectedOrders[i], line.Order(), "Order mismatch at index %d", i,
				)
			}
		})
	}
}

func BenchmarkLineFeed_NoLF(b *testing.B) {
	line := line.FromString(
		strings.Repeat("Hello World ", 100),
	)

	b.ReportAllocs()

	for b.Loop() {
		LineFeed(false, line)
	}
}

func BenchmarkLineFeed_SomeLF(b *testing.B) {
	line := line.FromString(
		strings.Repeat("Hello\nWorld\n", 100),
	)

	b.ReportAllocs()

	for b.Loop() {
		LineFeed(false, line)
	}
}
