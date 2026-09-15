package delta

import (
	"github.com/Rafael24595/go-reacterm-core/engine/helper/runes"
	"github.com/Rafael24595/go-reacterm-core/engine/model/offset"
)

// Delta defines a text replacement operation over a target offset range [Start, End).
type Delta struct {
	// Start is the starting offset of the range to be replaced.
	Start offset.Offset
	// End is the ending offset of the range to be replaced.
	End   offset.Offset
	// Text is the new text that will replace the specified range.
	Text  string
}

// New creates a new Delta mutation.
func New(start, end offset.Offset, text string) Delta {
	return Delta{
		Start: start,
		End:   end,
		Text:  text,
	}
}

// Measure returns the length of the delta's text payload in offset units.
func (d Delta) Measure() offset.Offset {
	return runes.MeasureOffset(d.Text)
}

// Apply executes the delta substitution against a rune buffer, returning a new modified slice.
func Apply(buffer []rune, delta *Delta) []rune {
	if delta == nil || delta.Start > delta.End {
		return buffer
	}

	size := offset.Offset(len(buffer))
	if delta.Start > size || delta.End > size {
		return buffer
	}

	deltaBuffer := runes.SanitizeRunes(
		[]rune(delta.Text),
	)

	runesSize := runes.MeasureOffsetRunes(deltaBuffer)

	tail := size - delta.End
	total := delta.Start + runesSize + tail

	newBuffer := make([]rune, total)

	copy(newBuffer[:delta.Start], buffer[:delta.Start])
	copy(newBuffer[delta.Start:], deltaBuffer)
	copy(newBuffer[delta.Start+runesSize:], buffer[delta.End:])

	return newBuffer
}
