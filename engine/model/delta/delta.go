package delta

import (
	"github.com/Rafael24595/go-reacterm-core/engine/helper/runes"
	"github.com/Rafael24595/go-reacterm-core/engine/model/offset"
)

type Delta struct {
	Start offset.Offset
	End   offset.Offset
	Text  string
}

func (d Delta) Measure() offset.Offset {
	return runes.Measureo(d.Text)
}

func Apply(buffer []rune, delta *Delta) []rune {
	if delta.Start > delta.End {
		return buffer
	}

	size := offset.Offset(len(buffer))
	if delta.Start > size || delta.End > size {
		return buffer
	}

	deltaBuffer := runes.SanitizeRunes(
		[]rune(delta.Text),
	)

	runesSize := runes.MeasureoRunes(deltaBuffer)

	tail := size - delta.End
	total := delta.Start + runesSize + tail

	newBuffer := make([]rune, total)

	copy(newBuffer[:delta.Start], buffer[:delta.Start])
	copy(newBuffer[delta.Start:], deltaBuffer)
	copy(newBuffer[delta.Start+runesSize:], buffer[delta.End:])

	return newBuffer
}
