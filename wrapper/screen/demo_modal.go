package wrapper_screen

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/primitive/modalmenu"
	"github.com/Rafael24595/go-reacterm-core/engine/model/input"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

func NewDemoModal() screen.Node {
	return modalmenu.New().
		SetName("modal - dolor").
		AddText(
			line.FromString("AD Lorem ipsum dolor sit amet"),
			line.Empty(),
		).
		AddOptions([]input.MenuOption{
			input.NewMenuOption("1", frag.FromString("Option_1")).WithHandler(NewDemoSelect),
			input.NewMenuOption("2", frag.FromString("Option_2")).WithHandler(NewDemoSelect),
			input.NewMenuOption("3", frag.FromString("Option_3")).WithHandler(NewDemoSelect),
			input.NewMenuOption("4", frag.FromString("Option_4")).WithHandler(NewDemoSelect),
		}...).
		ToNode()
}
