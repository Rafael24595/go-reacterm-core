package wrapper_screen

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/primitive/modalmenu"
	"github.com/Rafael24595/go-reacterm-core/engine/model/input"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
)

func NewDemoModal() screen.Node {
	return modalmenu.New().
		SetName("modal - dolor").
		AddText(
			*text.NewLine("AD Lorem ipsum dolor sit amet"),
			*text.EmptyLine(),
		).
		AddOptions([]input.MenuOption{
			input.NewMenuOption("1", *frag.New("Option_1"), NewDemoSelect),
			input.NewMenuOption("2", *frag.New("Option_2"), NewDemoSelect),
			input.NewMenuOption("3", *frag.New("Option_3"), NewDemoSelect),
			input.NewMenuOption("4", *frag.New("Option_4"), NewDemoSelect),
		}...).
		ToNode()
}
