package wrapper_screen

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/partial/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/partial/pipeline/header"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

func NewDemoHeader() pipeline.Transformer {
	return header.Transformer(
		pipeline.Before,
		line.SpecFrags(
			spec.AlignCenter(),
			frag.TextAtom(
				"Lorem ipsum dolor sit amet",
				atom.Upper,
			),
		),
		line.SpecFrags(
			spec.AlignCenter(),
			frag.TextAtom(
				"consectetur adipiscing",
				atom.Upper,
			),
		),
		line.SpecFrags(
			spec.AlignCenter(),
			frag.TextAtom(
				"-Server 00-",
				atom.Upper,
			),
		),
	)
}
