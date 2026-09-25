package marker

// TableSeparatorMeta defines the string glyphs used to frame and divide table elements.
type TableSeparatorMeta struct {
	// Top is the glyph used for the top border of the table.
	Top string
	// Bottom is the glyph used for the bottom border of the table.
	Bottom string
	// Center is the separator used between adjacent table cells.
	Center string
	// Left is the glyph or string used at the left edge of each table row.
	Left string
	// Right is the glyph or string used at the right edge of each table row.
	Right string
}

// DefaultTableSeparator defines a standard ASCII table framing style.
var DefaultTableSeparator = TableSeparatorMeta{
	Top:    "-",
	Bottom: "-",
	Center: " | ",
	Left:   "| ",
	Right:  " |",
}
