package pagination

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
)

var base_definition = screen.NewDefinition(
	map[key.Action]key.Descriptor{
		key.ActionPageUp:   {Code: []string{"⇞"}, Detail: "Prev page"},
		key.ActionPageDown: {Code: []string{"⇟"}, Detail: "Next page"},
	},
	[]key.Action{
		key.ActionPageUp,
		key.ActionPageDown,
	},
)

var definitions = map[pager.EngineCode]screen.Definition{
	pager.CodeEnginePaged:  pager_definition,
	pager.CodeEngineScroll: scroll_definition,
}

var keys = map[pager.EngineCode]struct {
	back key.Action
	next key.Action
}{
	pager.CodeEnginePaged:  {key.ActionArrowLeft, key.ActionArrowRight},
	pager.CodeEngineScroll: {key.ActionArrowUp, key.ActionArrowDown},
}

var labels = map[pager.EngineCode]string{
	pager.CodeEnginePaged:  "page",
	pager.CodeEngineScroll: "scroll",
}

var pager_definition = screen.NewDefinition(
	map[key.Action]key.Descriptor{
		key.ActionArrowLeft:  {Code: []string{"←"}, Detail: "Prev page"},
		key.ActionArrowRight: {Code: []string{"→"}, Detail: "Next page"},
	},
	[]key.Action{
		key.ActionArrowLeft,
		key.ActionArrowRight,
	},
)

var scroll_definition = screen.NewDefinition(
	map[key.Action]key.Descriptor{
		key.ActionArrowUp:   {Code: []string{"↑"}, Detail: "Scroll up"},
		key.ActionArrowDown: {Code: []string{"↓"}, Detail: "Scroll down"},
	},
	[]key.Action{
		key.ActionArrowUp,
		key.ActionArrowDown,
	},
)
