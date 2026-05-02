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

func Apply(buffer []rune, d *Delta) []rune {
	if d.Start > d.End {
		return buffer
	}

	size := offset.Offset(len(buffer))
	if d.Start > size || d.End > size {
		return buffer
	}

	runesSize := runes.Measureo(d.Text)

	tail := size - d.End
	total := d.Start + runesSize + tail

	newBuffer := make([]rune, total)

	copy(newBuffer[:d.Start], buffer[:d.Start])
	copy(newBuffer[d.Start:], []rune(d.Text))
	copy(newBuffer[d.Start+runesSize:], buffer[d.End:])

	return newBuffer
}
