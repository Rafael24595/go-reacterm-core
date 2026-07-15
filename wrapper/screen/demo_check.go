package wrapper_screen

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/partial/pipeline/header"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/primitive/checkmenu"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/runes"
	"github.com/Rafael24595/go-reacterm-core/engine/model/input"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

func NewDemoCheck() screen.Node {
	textTitle := "Sed facilisis, leo sit amet molestie congue, justo risus bibendum tortor"
	sizeTitle := runes.Measure(textTitle)

	title := []line.Line{
		*line.New(textTitle, spec.AlignLeft()),
		*line.New("-", spec.Fill(sizeTitle)),
	}

	options := []input.CheckOption{
		input.NewCheckOption("1", *frag.New("Check 1")),
		input.NewCheckOption("2", *frag.New("Check 2")),
		input.NewCheckOption("3", *frag.New("Check 3")),
		input.NewCheckOption("4", *frag.New("Check 4")),
	}

	node := checkmenu.New().
		Name("checkmenu - tortor").
		Limit(1).
		AddOptions(options...).
		ToNode()

	return header.Node(node, title...)
}
