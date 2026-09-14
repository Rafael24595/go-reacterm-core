package buffer

// Clipboard manages temporary in-memory storage for rune sequences.
type Clipboard struct {
	buffer []rune
}

// NewClipboard creates an empty Clipboard instance.
func NewClipboard() *Clipboard {
	return &Clipboard{}
}

// Size returns the number of runes stored in the clipboard.
func (c *Clipboard) Size() uint {
	return uint(len(c.buffer))
}

// Read returns a defensive copy of the stored runes to prevent external mutation.
func (c *Clipboard) Read() []rune {
	if len(c.buffer) == 0 {
		return nil
	}

	result := make([]rune, len(c.buffer))
	copy(result, c.buffer)
	return result
}

// Write replaces the clipboard content with a copy of the provided runes.
func (c *Clipboard) Write(rns []rune) *Clipboard {
	c.buffer = make([]rune, len(rns))
	copy(c.buffer, rns)
	return c
}

// Clear empties the clipboard content.
func (c *Clipboard) Clear() {
	c.buffer = nil
}
