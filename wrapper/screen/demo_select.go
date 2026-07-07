package wrapper_screen

import (
	"fmt"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/partial/pipeline/header"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/primitive/indexmenu"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/runes"
	"github.com/Rafael24595/go-reacterm-core/engine/model/input"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

func NewDemoSelect() screen.Node {
	textTitle := "Sed facilisis, leo sit amet molestie congue, justo risus bibendum tortor"
	sizeTitle := runes.Measure(textTitle)

	title := []text.Line{
		*text.NewLine(textTitle, spec.AlignLeft()),
		*text.NewLine("-", spec.Fill(sizeTitle)),
	}

	options := input.NewMenuOptions(
		input.NewMenuOption("opt_art", *text.NewFrag("[Prim] Option Article"), NewDemoArticle),
		input.NewMenuOption("opt_txt", *text.NewFrag("[Prim] Option TextArea"), NewDemoTextArea),
		input.NewMenuOption("opt_tbl", *text.NewFrag("[Prim] Option Table"), NewDemoTable),
		input.NewMenuOption("opt_mdl", *text.NewFrag("[Prim] Option Modal"), NewDemoModal),
		input.NewMenuOption("opt_chk", *text.NewFrag("[Prim] Option Check"), NewDemoCheck),
		input.NewMenuOption("opt_txi", *text.NewFrag("[Prim] Option TextInput"), NewDemoTextInput),
		input.NewMenuOption("opt_tlk", *text.NewFrag("[Prim] Option Talk"), NewDemoTalk),
		input.NewMenuOption("opt_clp", *text.NewFrag("[Prim] Option Clip"), NewDemoClip),
		input.NewMenuOption("opt_frm", *text.NewFrag("[Comp] Option Form"), NewDemoForm),
		input.NewMenuOption("opt_hsk", *text.NewFrag("[Demo] Option HStack"), NewDemoHStack),
	)

	optsSize := len(options)

	for i := range 30 {
		options = append(options,
			input.NewMenuOption(
				fmt.Sprintf("opt_%d", i),
				*text.NewFrag(fmt.Sprintf("Option %d", i+1+optsSize)),
				NewDemoTextArea,
			),
		)
	}

	node := indexmenu.New().
		SetName("indexmenu - tortor").
		SetMeta(marker.NumericIndex).
		AddOptions(options...).
		SetCursor(0).
		ToNode()

	return header.Node(node, title...)
}
