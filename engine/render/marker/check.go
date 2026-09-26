package marker

// CheckMeta defines the framing symbols and indicators used to render checkable options or checkboxes.
type CheckMeta struct {
	// Open is the left boundary glyph enclosing the check state (e.g. "[", "(").
	Open string
	// Close is the right boundary glyph enclosing the check state (e.g. "]", ")").
	Close string
	// Checked is the symbol displayed when the item is active or selected (e.g. "x", "✓", "*").
	Checked string
	// Unchecked is the symbol displayed when the item is inactive or unselected (e.g. " ", "-").
	Unchecked string
}

// BracketsCheck defines a standard ASCII bracket framing style (e.g. "[x]" / "[ ]").
var BracketsCheck = CheckMeta{
	Open:      "[",
	Close:     "]",
	Checked:   "x",
	Unchecked: " ",
}
