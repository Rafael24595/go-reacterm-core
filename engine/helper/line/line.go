package line

import (
	"github.com/Rafael24595/go-reacterm-core/engine/model/ascii"
	"github.com/Rafael24595/go-reacterm-core/engine/model/offset"
)

func DistanceFromLF(buffer []rune, from offset.Offset) offset.Offset {
	return from.Clamp(
		FindLineStart(buffer, from),
	)
}

func FindLineStart(buffer []rune, from offset.Offset) offset.Offset {
	if from == 0 {
		return 0
	}

	for i := from - 1; ; i-- {
		if buffer[i] == ascii.ENTER_LF {
			return i + 1
		}

		if i == 0 {
			break
		}
	}

	return 0
}

func FindLineEnd(buffer []rune, start offset.Offset) offset.Offset {
	i := start
	for i < offset.Offset(len(buffer)) && buffer[i] != ascii.ENTER_LF {
		i++
	}
	return i
}

func FindNextLineStart(buf []rune, from offset.Offset) (offset.Offset, bool) {
	for i := from; i < offset.Offset(len(buf)); i++ {
		if buf[i] == ascii.ENTER_LF {
			return i + 1, true
		}
	}
	return 0, false
}

func FindPrevLineStart(buf []rune, from offset.Offset) (offset.Offset, bool) {
	prevLineStart := FindLineStart(buf, from)
	if prevLineStart == 0 {
		return 0, false
	}
	return FindLineStart(buf, prevLineStart-1), true
}

func ClampToLine(buf []rune, lineStart, col offset.Offset) offset.Offset {
	end := FindLineEnd(buf, lineStart)
	lineLen := end - lineStart

	if col > lineLen {
		return end
	}

	return lineStart + col
}
