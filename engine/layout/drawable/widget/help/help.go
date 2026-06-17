package help

import (
	"fmt"
	"strings"

	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/drain"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

const Name = "help_unit"

type HelpUnit struct {
	loaded bool
	fields []key.Descriptor
	unit   drawable.Unit
}

func New(fields []key.Descriptor) *HelpUnit {
	return &HelpUnit{
		loaded: false,
		fields: fields,
	}
}

func UnitFromFields(fields []key.Descriptor) drawable.Unit {
	return New(fields).ToUnit()
}

func (u *HelpUnit) ToUnit() drawable.Unit {
	return drawable.NewBuilder().
		Name(Name).
		Boot(u.boot).
		Wipe(u.wipe).
		Draw(u.draw).
		ToUnit()
}

func (u *HelpUnit) boot() {
	u.loaded = true

	u.unit = makeUnit(u.fields)

	u.unit.Drawable.Boot()
}

func (u *HelpUnit) wipe() {
	if u.unit.Drawable.Wipe == nil {
		return
	}
	u.unit.Drawable.Wipe()
}

func (u *HelpUnit) draw(size winsize.Winsize) ([]text.Line, bool) {
	assert.True(u.loaded, drawable.MessageInitialized)

	return u.unit.Drawable.Draw(size)
}

func makeUnit(fields []key.Descriptor) drawable.Unit {
	fieldsLen := len(fields)
	if fieldsLen == 0 {
		return drain.UnitFromLines()
	}

	frags := make([]text.Fragment, fieldsLen)

	for i, field := range fields {
		code := strings.Join(field.Code, ", ")

		separator := ""
		if i < fieldsLen-1 {
			separator = " | "
		}

		frag := fmt.Sprintf("[%s] %s%s", code, field.Detail, separator)
		frags[i] = *text.NewFragment(frag).
			AddAtom(atom.Wrap)

	}

	return drain.UnitFromLines(
		*text.LineFromFragments(
			*text.NewFragment("--Help--"),
			*text.NewFragment("-").
				AddSpec(spec.Cover()),
		),
		*text.LineFromFragments(frags...),
		*text.NewLine("-", spec.Cover()),
	)
}
