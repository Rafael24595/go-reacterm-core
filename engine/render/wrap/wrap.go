package wrap

import (
	"sync"

	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
	"github.com/Rafael24595/go-reacterm-core/engine/render/wrap/layout"
	"github.com/Rafael24595/go-reacterm-core/engine/render/wrap/processor"
	"github.com/Rafael24595/go-reacterm-core/engine/render/wrap/splitter"
)

var once sync.Once
var wrapper = NewWrapper()

type Wrapper struct {
	processors []processor.Line
	splitter   splitter.Line
}

func NewWrapper(opts ...Option) Wrapper {
	return FromWrapper(
		DefaultWrapper(), opts...,
	)
}

func FromWrapper(cfg Wrapper, opts ...Option) Wrapper {
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}

func DefineWrapper(w Wrapper) bool {
	set := false
	once.Do(func() {
		wrapper = w
		set = true
	})
	return set
}

func (w Wrapper) normalize(order bool, lines ...line.Line) []layout.Line {
	buffer := make([]layout.Line, 0, len(lines))

	processed := w.proccess(order, lines...)

	for _, lne := range processed {
		words, frags := w.splitter(lne)

		layout := layout.NewLine(lne, words, frags)
		buffer = append(buffer, *layout)
	}

	return buffer
}

func (w Wrapper) proccess(order bool, lines ...line.Line) []line.Line {
	processed := lines
	for _, p := range w.processors {
		step := make([]line.Line, 0, len(processed)*2)
		for _, lne := range processed {
			step = append(
				step, p(order, lne)...,
			)
		}
		processed = step
	}

	return processed
}

func NormalizeLines(lines ...line.Line) []layout.Line {
	return wrapper.normalize(false, lines...)
}

func NormalizeLinesWithOrder(lines ...line.Line) []layout.Line {
	return wrapper.normalize(true, lines...)
}

func MaterializeEmpty(
	size winsize.Winsize,
	placeholder string,
	lines ...layout.Line,
) []layout.Line {
	for i, lne := range lines {
		if line.FragsMeasure(size.Cols, lne.Source) != 0 {
			continue
		}

		lastFrag := frag.Frag{}

		lneSize := lne.Source.Size()
		if lneSize > 0 {
			lastFrag = lne.Source.AtOrZero(lneSize - 1)
		}

		frag := frag.NewBuilder().
			AddText(placeholder).
			WithMeta(lastFrag).
			Frag()

		lines[i].Source = line.BuilderFromLine(lines[i].Source).
			PushFrags(frag).
			Line()

		lines[i].PushFrags(frag)
	}

	return lines
}

func Line(cols winsize.Cols, lne line.Line) []line.Line {
	return wrapLine(cols, lne, make([]line.Line, 0, 2))
}

func Lines(cols winsize.Cols, lines ...line.Line) []line.Line {
	result := make([]line.Line, 0, len(lines)*2)

	for _, line := range lines {
		result = wrapLine(cols, line, result)
	}

	return result
}

func wrapLine(cols winsize.Cols, line line.Line, dst []line.Line) []line.Line {
	normalized := NormalizeLines(line)

	for _, layout := range normalized {
		current := &layout

		for current != nil {
			head, rest := wrapOnce(cols, current)
			dst = append(dst, head.Line())
			current = rest
		}
	}

	return dst
}

func NextLine(cols winsize.Cols, lines []layout.Line) (*line.Line, []layout.Line) {
	builder, remain := NextBuilder(cols, lines)
	return builder.LinePtr(), remain
}

func NextBuilder(cols winsize.Cols, lines []layout.Line) (*line.Builder, []layout.Line) {
	if cols == 0 || len(lines) == 0 {
		return nil, make([]layout.Line, 0, len(lines)*2)
	}

	current := lines[0]
	remain := lines[1:]

	result, rest := wrapOnce(cols, &current)
	if rest != nil {
		remain = append([]layout.Line{*rest}, remain...)
	}

	return result, remain
}

func wrapOnce(cols winsize.Cols, lne *layout.Line) (*line.Builder, *layout.Line) {
	size := lne.Source.Size()

	cursor := line.NewBuilder(int(size)).
		WithMeta(lne.Source)

	remaining := cols
	currentWidth := winsize.Cols(0)

	wordIdx := uint(0)

	for ; wordIdx < lne.Size(); wordIdx++ {
		wordMeasure := lne.Measure(wordIdx, cols)

		if wordMeasure <= remaining {
			cursor.Text = layout.AppendFrags(
				cursor.Text, lne.FindFrags(wordIdx),
			)

			remaining = remaining.Sub(wordMeasure)
			currentWidth += wordMeasure

			continue
		}

		if shouldWrap(lne, wordIdx, currentWidth) {
			break
		}

		if ok := lne.SplitWord(
			wordIdx,
			cols,
			remaining,
		); ok {
			cursor.Text = layout.AppendFrags(
				cursor.Text, lne.FindFrags(wordIdx),
			)
		}

		wordIdx++

		break
	}

	if wordIdx >= lne.Size() {
		return cursor, nil
	}

	rest := layout.NewLine(
		lne.Source,
		lne.Words()[wordIdx:],
		lne.Frags(),
	)

	return cursor, rest
}

func shouldWrap(line *layout.Line, wordIdx uint, currentWidth winsize.Cols) bool {
	if line.HasAtom(wordIdx, atom.Break) {
		return false
	}

	return currentWidth > 0
}
