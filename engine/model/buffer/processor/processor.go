package processor

import (
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
)

// Value represents the underlying buffer of a text input.
type Value []rune

// Facade represents the display version of the buffer, which may differ from the underlying value (e.g., for masked inputs).
type Facade []rune

// Processor defines a function that receives an input buffer and returns the underlying buffer and its display facade.
type Processor func([]rune) (Value, Facade)

// Identity returns the buffer unmodified both as value and facade.
func Identity(buffer []rune) (Value, Facade) {
	return buffer, buffer
}

// Numeric filters non-numeric runes, allowing a single sign at the start and a single decimal separator (comma or dot).
func Numeric(buffer []rune) (Value, Facade) {
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

// Masked obscures non-whitespace characters with asterisks for password/secret inputs.
func Masked(buffer []rune) (Value, Facade) {
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

// Flatten transforms multiline text into single-line by replacing newlines with spaces.
func Flatten(buffer []rune) (Value, Facade) {
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

// Limit creates a decorator Processor that truncates both buffer and facade to max columns.
func Limit(
	limit winsize.Cols,
	processor Processor,
) Processor {
	if limit == 0 {
		return processor
	}

	return func(buffer []rune) (Value, Facade) {
		buffer, facade := processor(buffer)
		return trimBuffer(limit, facade, buffer)
	}
}

func trimBuffer(
	limit winsize.Cols,
	facade, buffer []rune,
) (Value, Facade) {
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
