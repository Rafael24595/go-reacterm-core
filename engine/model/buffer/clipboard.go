package buffer

type Clipboard struct {
	buffer []rune
}

func NewClipboard() *Clipboard {
	return &Clipboard{}
}

func (c *Clipboard) Size() uint {
	return uint(len(c.buffer))
}

func (c *Clipboard) Read() []rune {
	if len(c.buffer) == 0 {
		return nil
	}

	result := make([]rune, len(c.buffer))
	copy(result, c.buffer)
	return result
}

func (c *Clipboard) Write(rns []rune) *Clipboard {
	c.buffer = make([]rune, len(rns))
	copy(c.buffer, rns)
	return c
}

func (c *Clipboard) Clear() {
	c.buffer = nil
}
