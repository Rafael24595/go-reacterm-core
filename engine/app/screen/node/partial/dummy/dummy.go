package dummy

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/partial/template"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
)

const Name = "dummy"
const Tag = "system_dummy"

func ToNode() screen.Node {
	dummy := template.New().
		Name(Name).
		ViewModel(*viewmodel.New()).
		ToNode()

	dummy.Tags.Add(Tag)

	return dummy
}
