package buffer

import (
	"github.com/Rafael24595/go-reacterm-core/engine/helper/runes"
	"github.com/Rafael24595/go-reacterm-core/engine/model/buffer/processor"
	"github.com/Rafael24595/go-reacterm-core/engine/model/buffer/rule"
	"github.com/Rafael24595/go-reacterm-core/engine/model/delta"
	"github.com/Rafael24595/go-reacterm-core/engine/model/offset"
)

// RuneBuffer manages an in-memory sequence of runes, supporting transformation rules, display processors, and mutation tracking via incremental versions.
type RuneBuffer struct {
	buffer    []rune
	facade    []rune
	rules     []rule.Rule
	processor processor.Processor
	version   uint64
}

// NewRuneBuffer creates a new RuneBuffer initialized with default identity processor.
func NewRuneBuffer() *RuneBuffer {
	return &RuneBuffer{
		processor: processor.Identity,
	}
}

// WithRules appends auto-formatting rules to the buffer's execution pipeline.
func (b *RuneBuffer) WithRules(rules ...rule.Rule) *RuneBuffer {
	b.rules = append(b.rules, rules...)
	return b
}

// WithProcessor configures the display/value processor for the buffer.
func (b *RuneBuffer) WithProcessor(processor processor.Processor) *RuneBuffer {
	if processor != nil {
		b.processor = processor
	}
	return b
}

// Version returns the current mutation counter of the buffer.
func (b *RuneBuffer) Version() uint64 {
	return b.version
}

// Size returns the length of the underlying rune buffer.
func (b *RuneBuffer) Size() offset.Offset {
	return offset.Offset(len(b.buffer))
}

// Empty reports whether the buffer contains no runes.
func (b *RuneBuffer) Empty() bool {
	return len(b.buffer) == 0
}

// Buffer returns a defensive copy of the underlying stored runes.
func (b *RuneBuffer) Buffer() []rune {
	if len(b.buffer) == 0 {
		return nil
	}

	out := make([]rune, len(b.buffer))
	copy(out, b.buffer)
	return out
}

// Facade returns a defensive copy of the processed display runes.
func (b *RuneBuffer) Facade() []rune {
	if len(b.facade) == 0 {
		return nil
	}

	out := make([]rune, len(b.facade))
	copy(out, b.facade)
	return out
}

// Range returns a defensive copy of the runes within the specified offset boundaries.
func (b *RuneBuffer) Range(start, end offset.Offset) []rune {
	if end < start {
		return nil
	}

	buffLen := offset.Offset(len(b.buffer))
	if end > buffLen {
		end = buffLen
	}

	out := make([]rune, end-start)
	copy(out, b.buffer[start:end])
	return out
}

// Append inserts runes at the end of the current buffer.
func (b *RuneBuffer) Append(buffer []rune) *RuneBuffer {
	b.Replace(buffer, b.Size(), b.Size())
	return b
}

// Replace substitutes the runes between start and end offsets with the provided input buffer.
func (b *RuneBuffer) Replace(
	buffer []rune,
	start, end offset.Offset,
) ([]rune, []rune) {
	if end < start {
		return nil, nil
	}

	buffer = runes.SanitizeRunes(buffer)
	return b.commitReplace(buffer, start, end)
}

// ReplaceWithRules applies registered rules to the input buffer before committing the replacement.
func (b *RuneBuffer) ReplaceWithRules(
	buffer []rune,
	start, end offset.Offset,
) ([]rune, []rune) {
	if end < start {
		return nil, nil
	}

	buffer = runes.SanitizeRunes(buffer)
	buffer = b.applyRules(buffer, start, end, b.buffer)
	return b.commitReplace(buffer, start, end)
}

// Delete removes the range of runes between start and end offsets.
func (b *RuneBuffer) applyRules(
	buffer []rune,
	start, end offset.Offset,
	buff []rune,
) []rune {
	for _, rule := range b.rules {
		if text, ok := rule(buffer, start, end, buff); ok {
			return text
		}
	}
	return buffer
}

// Delete removes the range of runes between start and end offsets.
func (b *RuneBuffer) Delete(start, end offset.Offset) []rune {
	if end < start {
		return nil
	}

	_, deleted := b.Replace(nil, start, end)
	return deleted
}

func (b *RuneBuffer) commitReplace(
	buffer []rune,
	start, end offset.Offset,
) ([]rune, []rune) {
	buffLen := offset.Offset(len(b.buffer))

	if start > buffLen {
		start = buffLen
	}

	if end < start {
		end = start
	}

	if end > buffLen {
		end = buffLen
	}

	deleted := b.Range(start, end)

	rawBuffer := runes.Replace(b.buffer, buffer, start, end)
	newBuffer, newFacade := b.processor(rawBuffer)

	newBufferLen := offset.Offset(len(newBuffer))

	insertSize := offset.Offset(
		len(newBuffer) - (len(b.buffer) - len(deleted)),
	)

	var fixedInsert []rune
	if insertSize > 0 && start <= newBufferLen {
		endInsert := min(start+insertSize, newBufferLen)

		fixedInsert = make([]rune, endInsert-start)
		copy(fixedInsert, newBuffer[start:endInsert])
	}

	b.buffer = newBuffer
	b.facade = newFacade

	b.version += 1

	return fixedInsert, deleted
}

// ApplyDelta transforms the buffer state using a structural text delta.
func (b *RuneBuffer) ApplyDelta(d *delta.Delta) *RuneBuffer {
	newBuffer := delta.Apply(b.buffer, d)
	buffer, facade := b.processor(newBuffer)

	b.buffer = buffer
	b.facade = facade

	b.version += 1

	return b
}

// Clear resets the buffer and facade to an empty state, incrementing the version counter.
func (b *RuneBuffer) Clear() *RuneBuffer {
	b.buffer = nil
	b.facade = nil

	b.version += 1

	return b
}
