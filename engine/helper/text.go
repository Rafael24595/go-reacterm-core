package helper

import (
	"fmt"

	"github.com/Rafael24595/go-reacterm-core/engine/helper/runes"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
)

type Text struct {
	Data string
	Size winsize.Cols
}

func NewText(data string, size winsize.Cols) Text {
	return Text{
		Data: data,
		Size: size,
	}
}

func TextFromString(data string) Text {
	return NewText(data, runes.Measure(data))
}

func TextFromAny(data any) Text {
	return TextFromString(
		fmt.Sprintf("%v", data),
	)
}

func EmptyText() Text {
	return TextFromString("")
}

type Ellipsis struct {
	Data string
	Size winsize.Cols
}

func NewEllipsis(data string, size winsize.Cols) Ellipsis {
	return Ellipsis{
		Data: data,
		Size: size,
	}
}

func (e Ellipsis) measure() winsize.Cols {
	return runes.Measure(e.Data) * e.Size
}
