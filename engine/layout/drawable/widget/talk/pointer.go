package talk

import (
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

type pointerProvider func(cursor uint16, index uint16) ([]text.Fragment, []text.Fragment)

var providers = []pointerProvider{
	arrowProvider(),
	arrowProvider(marker.BlackRightPointingTriangle),
	gutterProvider(),
}

func FindPointer(cursor uint8) pointerProvider {
	if cursor >= uint8(len(providers)) {
		return providers[0]
	}
	return providers[cursor]
}

func NextPointer(cursor uint8) uint8 {
	return (cursor + 1) % uint8(len(providers))
}

func arrowProvider(arrow ...string) pointerProvider {
	pointer := marker.DefaultPromptText
	if len(arrow) > 0 {
		pointer = arrow[0]
	}

	defaultOwner := []text.Fragment{
		*text.NewFragment(marker.DefaultPaddingText).
			AddSpec(style.SpecRepeatRight(3)),
	}

	selectedOwner := []text.Fragment{
		*text.NewFragment(
			marker.DefaultPaddingText + pointer + marker.DefaultPaddingText,
		),
	}

	defaultMessage := []text.Fragment{
		*text.NewFragment(marker.DefaultPaddingText).
			AddSpec(style.SpecRepeatRight(5)),
	}

	return func(cursor, index uint16) ([]text.Fragment, []text.Fragment) {
		if index == cursor {
			return selectedOwner, defaultMessage
		}

		return defaultOwner, defaultMessage
	}
}

func gutterProvider() pointerProvider {
	defaultOwner := []text.Fragment{
		*text.NewFragment(marker.DefaultPaddingText).
			AddSpec(style.SpecRepeatRight(3)),
	}

	selectedOwner := []text.Fragment{
		*text.NewFragment(marker.DefaultMiddleGutterText),
		*text.NewFragment(marker.DefaultPaddingText).
			AddSpec(style.SpecRepeatRight(2)),
	}

	defaultMessage := []text.Fragment{
		*text.NewFragment(marker.DefaultPaddingText).
			AddSpec(style.SpecRepeatRight(5)),
	}

	selectedMessage := []text.Fragment{
		*text.NewFragment(marker.DefaultMiddleGutterText),
		*text.NewFragment(marker.DefaultPaddingText).
			AddSpec(style.SpecRepeatRight(4)),
	}

	return func(cursor, index uint16) ([]text.Fragment, []text.Fragment) {
		if index == cursor {
			return selectedOwner, selectedMessage
		}

		return defaultOwner, defaultMessage
	}
}
