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
		*line.FromFrags(
			*frag.New("Lorem ipsum dolor sit amet").
				AddSpec(spec.AlignCenter()).
				AddAtom(atom.Upper),
		),
		*line.FromFrags(
			*frag.New("consectetur adipiscing").
				AddSpec(spec.AlignCenter()).
				AddAtom(atom.Upper),
		),
		*line.FromFrags(
			*frag.New("-Server 00-").
				AddSpec(spec.AlignCenter()).
				AddAtom(atom.Upper),
		),
	)
}
