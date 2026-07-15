package padding

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/config/padding/cols"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/transform"
	"github.com/Rafael24595/go-reacterm-core/engine/model/hint"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

type colPositioner func(winsize.Cols) (winsize.Cols, winsize.Cols)

var colPositionerMap = map[style.HorizontalPosition]colPositioner{
	style.Left:   colToLeft,
	style.Right:  colToRight,
	style.Center: colToCenter,
}

func Cols(hint hint.Size[winsize.Cols], opts ...cols.Option) transform.Transformer {
	return func(size winsize.Winsize, lines []line.Line) []line.Line {
		newLines := make([]line.Line, len(lines))
		fixedMin := hint.Min(size.Cols)

		for i := range lines {
			remaining := fixedMin.Sub(
				frag.Measure(size.Cols, lines[i].Text...),
			)

			if remaining == 0 {
				newLines[i] = lines[i]
				continue
			}

			newLines[i] = AddColsPadding(remaining, lines[i], opts...)
		}

		return newLines
	}
}

func colToLeft(remaining winsize.Cols) (winsize.Cols, winsize.Cols) {
	return 0, remaining
}

func colToRight(remaining winsize.Cols) (winsize.Cols, winsize.Cols) {
	return remaining, 0
}

func colToCenter(remaining winsize.Cols) (winsize.Cols, winsize.Cols) {
	paddingL := remaining / 2
	paddingR := remaining.Sub(paddingL)
	return paddingL, paddingR
}

func AddColsPadding(
	size winsize.Cols,
	lne line.Line,
	opts ...cols.Option,
) line.Line {
	cfg := cols.ResolveConfig(opts...)

	frg := cfg.Provider(size, lne)

	positioner, ok := colPositionerMap[cfg.Position]
	if !ok {
		assert.Unreachable("undefined vertical position '%d'", cfg.Position)
		positioner = colToLeft
	}

	paddingL, paddingR := positioner(size)

	frags := make([]frag.Frag, 0, 3)

	if paddingL > 0 {
		frag := frg.Clone().
			AddSpec(spec.ExtendRight(paddingL))
		frags = append(frags, *frag)
	}

	frags = append(frags, lne.Text...)

	if paddingR > 0 {
		frag := frg.Clone().
			AddSpec(spec.ExtendRight(paddingR))
		frags = append(frags, *frag)
	}

	return *line.FromMeta(&lne).
		PushFrags(frags...)
}
