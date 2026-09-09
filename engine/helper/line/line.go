package line

import (
	"github.com/Rafael24595/go-reacterm-core/engine/model/ascii"
	"github.com/Rafael24595/go-reacterm-core/engine/model/offset"
)

// DistanceFromLF calculates the offset distance from the start of the current line to the given position.
func DistanceFromLF(buffer []rune, from offset.Offset) offset.Offset {
	return from.Sub(
		FindLineStart(buffer, from),
	)
}

// FindLineStart locates the offset corresponding to the beginning of the line containing 'from'.
func FindLineStart(buffer []rune, from offset.Offset) offset.Offset {
	bufLen := offset.Offset(len(buffer))
	if bufLen == 0 || from == 0 {
		return 0
	}

	if from > bufLen {
		from = bufLen
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

// FindLineEnd finds the offset of the newline character or end of buffer starting from 'start'.
func FindLineEnd(buffer []rune, start offset.Offset) offset.Offset {
	bufLen := offset.Offset(len(buffer))

	i := start
	for i < bufLen && buffer[i] != ascii.ENTER_LF {
		i++
	}

	return i
}

// FindNextLineStart finds the offset of the first character following the next LF byte.
func FindNextLineStart(buf []rune, from offset.Offset) (offset.Offset, bool) {
	bufLen := offset.Offset(len(buf))

	for i := from; i < bufLen; i++ {
		if buf[i] == ascii.ENTER_LF {
			return i + 1, true
		}
	}
	
	return 0, false
}

// FindPrevLineStart finds the starting offset of the line prior to the current line at 'from'.
func FindPrevLineStart(buf []rune, from offset.Offset) (offset.Offset, bool) {
	prevLineStart := FindLineStart(buf, from)
	if prevLineStart == 0 {
		return 0, false
	}
	return FindLineStart(buf, prevLineStart-1), true
}

// ClampToLine constrains a column offset to the actual available length of the target line.
func ClampToLine(buf []rune, lineStart, col offset.Offset) offset.Offset {
	end := FindLineEnd(buf, lineStart)
	lineLen := end - lineStart

	if col > lineLen {
		return end
	}

	return lineStart + col
}
