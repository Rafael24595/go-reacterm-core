package buffer

import (
	"github.com/Rafael24595/go-reacterm-core/engine/helper/runes"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/text"
	"github.com/Rafael24595/go-reacterm-core/engine/model/delta"
	"github.com/Rafael24595/go-reacterm-core/engine/model/offset"
)

type RuneBuffer struct {
	buffer      []rune
	facade      []rune
	transformer text.TextTransformer
	handler     RuneHandler
}

func NewRuneBuffer() *RuneBuffer {
	return &RuneBuffer{
		buffer:      make([]rune, 0),
		facade:      make([]rune, 0),
		transformer: text.VoidTextTransformer,
		handler:     voidRuneHandler,
	}
}

func (b *RuneBuffer) Transformer(transformer text.TextTransformer) *RuneBuffer {
	b.transformer = transformer
	return b
}

func (b *RuneBuffer) Handler(handler RuneHandler) *RuneBuffer {
	b.handler = handler
	return b
}

func (b *RuneBuffer) Size() offset.Offset {
	return offset.Offset(len(b.buffer))
}

func (b *RuneBuffer) Empty() bool {
	return len(b.buffer) == 0
}

func (b *RuneBuffer) Buffer() []rune {
	return b.buffer
}

func (b *RuneBuffer) Facade() []rune {
	return b.facade
}

func (b *RuneBuffer) Range(start offset.Offset, end offset.Offset) []rune {
	if end < start {
		return make([]rune, 0)
	}

	return b.buffer[start:end]
}

func (b *RuneBuffer) Append(rns []rune) *RuneBuffer {
	b.Replace(rns, b.Size(), b.Size())
	return b
}

func (b *RuneBuffer) TransformAndReplace(buffer []rune, start offset.Offset, end offset.Offset) ([]rune, []rune) {
	if end < start {
		zero := make([]rune, 0)
		return zero, zero
	}

	insert := b.transformer.Apply(buffer, start, end, b.buffer)
	return b.applyChange(insert, start, end)
}

func (b *RuneBuffer) Replace(rns []rune, start offset.Offset, end offset.Offset) ([]rune, []rune) {
	if end < start {
		zero := make([]rune, 0)
		return zero, zero
	}

	return b.applyChange(rns, start, end)
}

func (b *RuneBuffer) Delete(start offset.Offset, end offset.Offset) []rune {
	if end < start {
		return make([]rune, 0)
	}

	rns := make([]rune, 0)
	_, deleted := b.Replace(rns, start, end)
	return deleted
}

func (b *RuneBuffer) applyChange(insert []rune, start, end offset.Offset) ([]rune, []rune) {
	end = min(end, offset.Offset(len(b.buffer)))

	deleted := b.Range(start, end)

	rawBuffer := runes.AppendRange(b.buffer, insert, start, end)
	newBuffer, newFacade := b.handler(rawBuffer)

	newBufferLen := offset.Offset(len(newBuffer))

	insertSize := offset.Offset(
		len(newBuffer) - (len(b.buffer) - len(deleted)),
	)

	fixedInsert := make([]rune, 0)
	if insertSize > 0 {
		endInsert := min(start+insertSize, newBufferLen)
		fixedInsert = newBuffer[start:endInsert]
	}

	b.buffer = newBuffer
	b.facade = newFacade

	return fixedInsert, deleted
}

func (b *RuneBuffer) ApplyDelta(d *delta.Delta) *RuneBuffer {
	newBuffer := delta.Apply(b.buffer, d)
	buffer, facade := b.handler(newBuffer)

	b.buffer = buffer
	b.facade = facade

	return b
}
