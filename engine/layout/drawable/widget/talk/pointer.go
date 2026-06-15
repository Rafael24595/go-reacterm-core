package talk

import (
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

type PointerProvider func(cursor uint16, index uint16) ([]text.Fragment, []text.Fragment)

var providers = []PointerProvider{
	arrowProvider(),
	arrowProvider(marker.U25B6),
	gutterProvider(),
}

var NoneProvider = arrowProvider(' ')

func FindPointer(cursor uint8) PointerProvider {
	if cursor >= uint8(len(providers)) {
		return providers[0]
	}
	return providers[cursor]
}

func NextPointer(cursor uint8) uint8 {
	return (cursor + 1) % uint8(len(providers))
}

func arrowProvider(arrow ...rune) PointerProvider {
	pointer := marker.DefaultPromptText
	if len(arrow) > 0 {
		pointer = string(arrow[0])
	}

	defaultOwner := []text.Fragment{
		*text.NewFragment(marker.DefaultPaddingText).
			AddSpec(spec.ExtendRight(3)),
	}

	selectedOwner := []text.Fragment{
		*text.NewFragment(
			marker.DefaultPaddingText + pointer + marker.DefaultPaddingText,
		),
	}

	defaultMessage := []text.Fragment{
		*text.NewFragment(marker.DefaultPaddingText).
			AddSpec(spec.ExtendRight(5)),
	}

	return func(cursor, index uint16) ([]text.Fragment, []text.Fragment) {
		if index == cursor {
			return selectedOwner, defaultMessage
		}

		return defaultOwner, defaultMessage
	}
}

func gutterProvider() PointerProvider {
	defaultOwner := []text.Fragment{
		*text.NewFragment(marker.DefaultPaddingText).
			AddSpec(spec.ExtendRight(3)),
	}

	selectedOwner := []text.Fragment{
		*text.NewFragment(marker.DefaultMiddleGutterText),
		*text.NewFragment(marker.DefaultPaddingText).
			AddSpec(spec.ExtendRight(2)),
	}

	defaultMessage := []text.Fragment{
		*text.NewFragment(marker.DefaultPaddingText).
			AddSpec(spec.ExtendRight(5)),
	}

	selectedMessage := []text.Fragment{
		*text.NewFragment(marker.DefaultMiddleGutterText),
		*text.NewFragment(marker.DefaultPaddingText).
			AddSpec(spec.ExtendRight(4)),
	}

	return func(cursor, index uint16) ([]text.Fragment, []text.Fragment) {
		if index == cursor {
			return selectedOwner, selectedMessage
		}

		return defaultOwner, defaultMessage
	}
}
