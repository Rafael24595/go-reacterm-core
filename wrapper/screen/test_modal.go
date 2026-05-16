package wrapper_screen

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/primitive/modalmenu"
	"github.com/Rafael24595/go-reacterm-core/engine/model/input"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

func NewTestModal() screen.Node {
	return modalmenu.New().
		SetName("modal - dolor").
		AddText(
			*text.NewLine("AD Lorem ipsum dolor sit amet"),
			*text.EmptyLine(),
		).
		AddOptions([]input.MenuOption{
			input.NewMenuOption("1", *text.NewFragment("Option_1"), NewTestSelect),
			input.NewMenuOption("2", *text.NewFragment("Option_2"), NewTestSelect),
			input.NewMenuOption("3", *text.NewFragment("Option_3"), NewTestSelect),
			input.NewMenuOption("4", *text.NewFragment("Option_4"), NewTestSelect),
		}...).
		ToNode()
}
