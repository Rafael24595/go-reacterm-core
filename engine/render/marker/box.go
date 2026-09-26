package marker

// BoxSeparatorMeta defines the boundary glyphs and internal filler character used to frame rectangular containers or dialog boxes.
type BoxSeparatorMeta struct {
	// Top is the horizontal line glyph used for the top border (e.g. "-", "─").
	Top string
	// Bottom is the horizontal line glyph used for the bottom border (e.g. "-", "─").
	Bottom string
	// Left is the vertical line glyph used for the left border (e.g. "|", "│").
	Left string
	// Right is the vertical line glyph used for the right border (e.g. "|", "│").
	Right string
	// Space is the character glyph used for internal padding or blank cell fills (e.g. " ").
	Space string
}

// DefaultBoxSeparator defines a standard ASCII box framing style.
var DefaultBoxSeparator = BoxSeparatorMeta{
	Top:    "-",
	Bottom: "-",
	Left:   "|",
	Right:  "|",
	Space:  " ",
}
