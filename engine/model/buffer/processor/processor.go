package processor

import (
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
)

type Value []rune
type Facade []rune

type Processor func([]rune) (Value, Facade)

func Identity(buffer []rune) (Value, Facade) {
	return buffer, buffer
}

func Number(buffer []rune) (Value, Facade) {
	fixedBuffer := make([]rune, 0, len(buffer))
	hasSeparator := false

	for _, v := range buffer {
		if v >= '0' && v <= '9' {
			fixedBuffer = append(fixedBuffer, v)
			continue
		}

		if (v == ',' || v == '.') && !hasSeparator {
			fixedBuffer = append(fixedBuffer, v)
			hasSeparator = true
			continue
		}

		if v == '-' && len(fixedBuffer) == 0 {
			fixedBuffer = append(fixedBuffer, v)
		}
	}

	if len(fixedBuffer) == len(buffer) {
		return buffer, buffer
	}

	return fixedBuffer, fixedBuffer
}

func Hidden(buffer []rune) (Value, Facade) {
	fixedBuffer := make([]rune, len(buffer))

	for i, r := range buffer {
		switch r {
		case '\n', '\t', ' ':
			fixedBuffer[i] = r
		default:
			fixedBuffer[i] = '*'
		}
	}

	return buffer, fixedBuffer
}

func Inline(buffer []rune) (Value, Facade) {
	fixedBuffer := make([]rune, len(buffer))

	for i, r := range buffer {
		switch r {
		case '\n':
			fixedBuffer[i] = ' '
		default:
			fixedBuffer[i] = buffer[i]
		}
	}

	return fixedBuffer, fixedBuffer
}

func Limit(limit winsize.Cols, processor Processor) Processor {
	if limit == 0 {
		return processor
	}

	return func(buffer []rune) (Value, Facade) {
		buffer, facade := processor(buffer)
		return trimBuffer(limit, facade, buffer)
	}
}

func trimBuffer(limit winsize.Cols, facade, buffer []rune) (Value, Facade) {
	if limit == 0 {
		return buffer, facade
	}

	if len(buffer) > int(limit) {
		buffer = buffer[:limit]
	}

	if len(facade) > int(limit) {
		facade = facade[:limit]
	}

	return buffer, facade
}
