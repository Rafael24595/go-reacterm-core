package layout

import (
	"slices"

	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

// TODO: Review documentation.

const errf_word_out_of_range = "index out of words range [%d] with length %d"
const errf_frag_out_of_range = "index out of frags range [%d] with length %d"

// Line is a mutable layout representation used during wrapping.
//
// A Line owns its words and fragments and may mutate them in place.
// Callers that need an independent snapshot must use clone/Clones.
type Line struct {
	Source line.Line
	words  []Word
	frags  []Frag
}

func NewLine(source line.Line, words []Word, frags []Frag) *Line {
	return &Line{
		Source: source,
		words:  words,
		frags:  frags,
	}
}

func (l *Line) Size() uint {
	return uint(len(l.words))
}

func (l *Line) Words() []Word {
	clone := make([]Word, len(l.words))
	copy(clone, l.words)
	return clone
}

func (l *Line) Frags() []Frag {
	clone := make([]Frag, len(l.frags))
	copy(clone, l.frags)
	return clone
}

func (l *Line) FindFrags(idx uint) []Frag {
	if idx >= uint(len(l.words)) {
		assert.Unreachable(errf_word_out_of_range, idx, len(l.words))
		return make([]Frag, 0)
	}

	word := l.words[idx]
	return l.frags[word.start:word.end]
}

func (l *Line) PushFrags(frags ...frag.Frag) *Line {
	lenFrags := len(l.frags)

	word := New(
		uint32(lenFrags),
		uint32(lenFrags+len(frags)),
	)

	l.words = append(l.words, *word)
	l.frags = AppendFromFrags(l.frags, frags)

	return l
}

// SliceFromWord removes the first idx words from the line.
//
// The operation mutates Line and retains the existing backing arrays.
func (l *Line) SliceFromWord(idx uint) *Line {
	if idx >= uint(len(l.words)) {
		assert.Unreachable(errf_word_out_of_range, idx, len(l.words))
		return l
	}

	l.words = l.words[idx:]
	return l
}

func (l *Line) SplitWord(
	wordIdx uint,
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
			wordIdx, fragIdx, remaining,
		)

		return true
	}

	return true
}

func (l *Line) splitFrag(
	wordIdx uint,
	fragIdx uint32,
	cols winsize.Cols,
) {
	if fragIdx >= uint32(len(l.frags)) {
		assert.Unreachable(errf_frag_out_of_range, fragIdx, len(l.words))
		return
	}

	if wordIdx >= uint(len(l.words)) {
		assert.Unreachable(errf_word_out_of_range, fragIdx, len(l.words))
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
		l.frags, int(nextIdx), *right,
	)

	oldEnd := l.words[wordIdx].end

	l.words[wordIdx].end = nextIdx
	l.words[wordIdx].measured = false
	l.words[wordIdx].measure = 0

	newWord := New(
		l.words[wordIdx].end,
		oldEnd+1,
	)

	l.words = slices.Insert(
		l.words, int(wordIdx+1), *newWord,
	)

	for i := wordIdx + 2; i < l.Size(); i++ {
		l.words[i].start++
		l.words[i].end++
	}
}

func (l *Line) HasAtom(idx uint, atm atom.Atom) bool {
	if idx >= uint(len(l.words)) {
		assert.Unreachable(errf_word_out_of_range, idx, len(l.words))
		return false
	}

	for _, v := range l.FindFrags(idx) {
		if v.Base.Atom().HasAny(atm) {
			return true
		}
	}
	return false
}

func (l *Line) Measure(idx uint, cols winsize.Cols) winsize.Cols {
	return l.measureWith(idx, cols, fragMeasure)
}

func (l *Line) measureWith(
	idx uint,
	cols winsize.Cols,
	resolver measureResolver,
) winsize.Cols {
	if idx >= uint(len(l.words)) {
		assert.Unreachable(errf_word_out_of_range, idx, len(l.words))
		return 0
	}

	word := &l.words[idx]

	if word.measured && word.cols == cols {
		return word.measure
	}

	measure := resolver(
		cols, l.FindFrags(idx)...,
	)

	word.cols = cols
	word.measure = measure
	word.measured = true

	return measure
}

func (l *Line) clone() *Line {
	newLine := NewLine(l.Source, l.words, l.frags)

	newLine.words = slices.Clone(l.words)
	newLine.frags = slices.Clone(l.frags)

	return newLine
}

func Clones(lines ...Line) []Line {
	clones := make([]Line, len(lines))
	for i, v := range lines {
		clones[i] = *v.clone()
	}
	return clones
}
