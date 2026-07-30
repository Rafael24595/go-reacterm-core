package wrap

import (
	"testing"

	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

func BenchmarkRender_FirstFrame(b *testing.B) {
	lne := benchmarkLine(20_000)

	b.ReportAllocs()

	for b.Loop() {
		NormalizeLines(lne)
	}
}

func BenchmarkRender_SecondFrame(b *testing.B) {
	lne := benchmarkLine(20_000)

	NormalizeLines(lne)

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		NormalizeLines(lne)
	}
}

func BenchmarkRender_InsertRune(b *testing.B) {
	text := benchmarkText(20_000)

	b.ReportAllocs()

	for b.Loop() {
		modified := text + "."

		lne := line.FromFrags(
			frag.FromStrings(modified)...,
		)

		NormalizeLines(lne)
	}
}

func BenchmarkRender_StyleChange(b *testing.B) {
	text := benchmarkText(20_000)

	b.ReportAllocs()

	for b.Loop() {
		lne := line.NewBuilder().
			PushText(text).
			AddSpec(
				spec.Fill(80),
			).
			Line()

		NormalizeLines(lne)
	}
}
