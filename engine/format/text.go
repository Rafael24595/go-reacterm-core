package format

import (
	"fmt"

	"github.com/Rafael24595/go-reacterm-core/engine/helper/runes"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
)

// Text wraps a raw string alongside its measured visual width in columns.
type Text struct {
	// Data holds the raw string content.
	Data string
	// Size represents the visual width of the string in columns.
	Size winsize.Cols
}

// NewText creates a Text instance with an explicit string value and column size.
func NewText(data string, size winsize.Cols) Text {
	return Text{
		Data: data,
		Size: size,
	}
}

// TextFromString creates a Text instance measuring the visual width automatically.
func TextFromString(data string) Text {
	return NewText(data, runes.MeasureCols(data))
}

// TextFromAny converts any value into its string representation and measures its width.
func TextFromAny(data any) Text {
	return TextFromString(
		fmt.Sprintf("%v", data),
	)
}

// EmptyText constructs an empty Text value with 0 column size.
func EmptyText() Text {
	return TextFromString("")
}

// IsEmpty returns true if the inner data string is empty or width is zero.
func (t Text) IsEmpty() bool {
	return t.Data == "" || t.Size == 0
}
