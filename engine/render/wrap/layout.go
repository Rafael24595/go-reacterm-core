package wrap

import (
	"slices"

	assert "github.com/Rafael24595/go-assert/assert/runtime"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

type LayoutLine struct {
	Source *text.Line
	words  []word
	frags  []wordFrag
}

func NewLayoutLine(source *text.Line, words []word, frags []wordFrag) *LayoutLine {
	return &LayoutLine{
		Source: source,
		words:  words,
		frags:  frags,
	}
}

func (l *LayoutLine) findFrags(idx int) []wordFrag {
	if idx >= len(l.words) {
		assert.Unreachable(
			"index out of words range [%d] with length %d", idx, len(l.words),
		)
		return make([]wordFrag, 0)
	}

	word := l.words[idx]
	return l.frags[word.start:word.end]
}

func (l *LayoutLine) pushFrags(frags ...text.Frag) *LayoutLine {
	lenFrags := len(l.frags)
	word := word{
		start: uint32(lenFrags),
		end:   uint32(lenFrags + len(frags)),
	}

	l.words = append(l.words, word)
	l.frags = append(l.frags,
		toWordFrag(frags...)...,
	)

	return l
}

func (l *LayoutLine) splitWord(
	wordIdx int,
	cols winsize.Cols,
	remaining winsize.Cols,
) bool {
	if cols == 0 || remaining == 0 {
		return false
	}

	current := &l.words[wordIdx]

	for fragIdx := current.start; fragIdx < current.end; fragIdx++ {
		frag := &l.frags[fragIdx]

		size := frag.Measure(cols)

		if size <= remaining {
			remaining = remaining.Sub(size)
			continue
		}

		l.splitFrag(
			wordIdx, int(fragIdx), remaining,
		)

		return true
	}

	return true
}

func (l *LayoutLine) splitFrag(
	wordIdx int,
	fragIdx int,
	cols winsize.Cols,
) {
	if fragIdx >= len(l.frags) {
		assert.Unreachable(
			"index out of frags range [%d] with length %d",
			fragIdx,
			len(l.words),
		)

		return
	}

	if wordIdx >= len(l.words) {
		assert.Unreachable(
			"index out of words range [%d] with length %d",
			fragIdx,
			len(l.words),
		)

		return
	}

	frag := &l.frags[fragIdx]

	left, right := splitFragAt(frag, cols)
	if right == nil {
		return
	}

	l.frags[fragIdx] = *left

	nextIdx := fragIdx + 1
	l.frags = slices.Insert(
		l.frags, nextIdx, *right,
	)

	oldEnd := l.words[wordIdx].end

	l.words[wordIdx].end = uint32(nextIdx)
	l.words[wordIdx].measured = false
	l.words[wordIdx].measure = 0

	newWord := word{
		start: l.words[wordIdx].end,
		end:   oldEnd + 1,
	}

	l.words = slices.Insert(
		l.words, wordIdx+1, newWord,
	)

	for i := wordIdx + 2; i < len(l.words); i++ {
		l.words[i].start++
		l.words[i].end++
	}
}

func (l *LayoutLine) hasAtom(idx int, atm atom.Atom) bool {
	if idx >= len(l.words) {
		assert.Unreachable(
			"index out of words range [%d] with length %d", idx, len(l.words),
		)
		return false
	}

	for _, v := range l.findFrags(idx) {
		if v.Base.Atom.HasAny(atm) {
			return true
		}
	}
	return false
}

func (l *LayoutLine) measure(idx int, cols winsize.Cols) winsize.Cols {
	return l.measureWith(idx, cols, fragMeasure)
}

func (l *LayoutLine) measureWith(
	idx int,
	cols winsize.Cols,
	resolver measureResolver,
) winsize.Cols {
	if idx >= len(l.words) {
		assert.Unreachable(
			"index out of words range [%d] with length %d", idx, len(l.words),
		)
		return 0
	}

	word := &l.words[idx]

	if word.measured && word.cols == cols {
		return word.measure
	}

	measure := resolver(
		cols, l.findFrags(idx)...,
	)

	word.cols = cols
	word.measure = measure
	word.measured = true

	return measure
}

func (l *LayoutLine) clone() *LayoutLine {
	newLine := NewLayoutLine(l.Source, l.words, l.frags)

	newLine.words = slices.Clone(l.words)
	newLine.frags = slices.Clone(l.frags)

	return newLine
}

func CloneLayoutLines(lines ...LayoutLine) []LayoutLine {
	clones := make([]LayoutLine, len(lines))
	for i, v := range lines {
		clones[i] = *v.clone()
	}
	return clones
}
