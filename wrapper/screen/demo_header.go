package wrapper_screen

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/partial/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/partial/pipeline/header"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
)

func NewDemoHeader() pipeline.Transformer {
	lines := text.ApplyLineSpec(
		spec.AlignCenter(),
		*text.LineFromFrags(
			*frag.New("Lorem ipsum dolor sit amet").AddAtom(atom.Upper),
		),
		*text.LineFromFrags(
			*frag.New("consectetur adipiscing").AddAtom(atom.Upper),
		),
		*text.LineFromFrags(
			*frag.New("-Server 00-").AddAtom(atom.Upper),
		),
	)

	return header.Transformer(pipeline.Before, lines...)
}
