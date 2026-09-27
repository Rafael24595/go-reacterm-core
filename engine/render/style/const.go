package style

// Direction specifies the layout orientation (Horizontal or Vertical).
type Direction uint8

const (
	// Horizontal places elements sequentially along the X-axis (left to right).
	Horizontal Direction = iota
	// Vertical places elements sequentially along the Y-axis (top to bottom).
	Vertical
)

// VerticalPosition defines vertical alignment within a container.
type VerticalPosition uint8

const (
	// Top aligns elements to the upper boundary of the container.
	Top VerticalPosition = iota
	// Middle aligns elements centrally along the vertical axis.
	Middle
	// Bottom aligns elements to the lower boundary of the container.
	Bottom
)

// HorizontalPosition defines horizontal alignment within a container.
type HorizontalPosition uint8

const (
	// Left aligns elements to the left boundary of the container.
	Left HorizontalPosition = iota
	// Center aligns elements centrally along the horizontal axis.
	Center
	// Right aligns elements to the right boundary of the container.
	Right
)

// Justify specifies content distribution along the main axis (similar to CSS Flexbox).
type Justify uint8

const (
	// JustifyStart packs elements toward the start line of the main axis.
	JustifyStart Justify = iota
	// JustifyEnd packs elements toward the end line of the main axis.
	JustifyEnd
	// JustifyCenter packs elements along the center of the main axis.
	JustifyCenter
	// JustifyBetween distributes elements evenly; first element is at the start, last at the end.
	JustifyBetween
	// JustifyAround distributes elements evenly with equal space around each element.
	JustifyAround
	// JustifyEvenly distributes elements so that the space between any two items is equal.
	JustifyEvenly
)
