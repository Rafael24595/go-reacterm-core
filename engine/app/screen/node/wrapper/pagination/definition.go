package pagination

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager/action"
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

var definitions = map[action.Kind]screen.Definition{
	action.KindPaged:  pager_definition,
	action.KindScroll: scroll_definition,
}

var keys = map[action.Kind]struct {
	back key.Action
	next key.Action
}{
	action.KindPaged:  {key.ActionArrowLeft, key.ActionArrowRight},
	action.KindScroll: {key.ActionArrowUp, key.ActionArrowDown},
}

var labels = map[action.Kind]string{
	action.KindPaged:  "page",
	action.KindScroll: "scroll",
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
