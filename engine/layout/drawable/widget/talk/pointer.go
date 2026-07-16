package talk

import (
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
)

type PointerProvider func(cursor uint16, index uint16) ([]frag.Frag, []frag.Frag)

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

	defaultOwner := []frag.Frag{
		frag.TextSpec(
			marker.DefaultPaddingText, spec.ExtendRight(3),
		),
	}

	selectedOwner := []frag.Frag{
		frag.FromString(
			marker.DefaultPaddingText + pointer + marker.DefaultPaddingText,
		),
	}

	defaultMessage := []frag.Frag{
		frag.TextSpec(
			marker.DefaultPaddingText, spec.ExtendRight(5),
		),
	}

	return func(cursor, index uint16) ([]frag.Frag, []frag.Frag) {
		if index == cursor {
			return selectedOwner, defaultMessage
		}

		return defaultOwner, defaultMessage
	}
}

func gutterProvider() PointerProvider {
	defaultOwner := []frag.Frag{
		frag.TextSpec(
			marker.DefaultPaddingText, spec.ExtendRight(3),
		),
	}

	selectedOwner := []frag.Frag{
		frag.FromString(marker.U2503_Text),
		frag.TextSpec(
			marker.DefaultPaddingText, spec.ExtendRight(2),
		),
	}

	defaultMessage := []frag.Frag{
		frag.TextSpec(
			marker.DefaultPaddingText, spec.ExtendRight(5),
		),
	}

	selectedMessage := []frag.Frag{
		frag.FromString(marker.U2503_Text),
		frag.TextSpec(
			marker.DefaultPaddingText, spec.ExtendRight(4),
		),
	}

	return func(cursor, index uint16) ([]frag.Frag, []frag.Frag) {
		if index == cursor {
			return selectedOwner, selectedMessage
		}

		return defaultOwner, defaultMessage
	}
}
