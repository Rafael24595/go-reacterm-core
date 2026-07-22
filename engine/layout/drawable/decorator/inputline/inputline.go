package inputline

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/transform/drain"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"

	drawable_drain "github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/drain"
)

const Name = "input_line_unit"

type InputLineUnit struct {
	loaded bool
	status bool
	prompt string
	unit   drawable.Unit
}

func New(unit drawable.Unit) *InputLineUnit {
	return &InputLineUnit{
		loaded: false,
		status: true,
		prompt: marker.DefaultPromptText,
		unit:   unit,
	}
}

func Wrap(unit drawable.Unit) drawable.Unit {
	return New(unit).ToUnit()
}

func FromString(text string) drawable.Unit {
	return Wrap(
		drawable_drain.UnitFromString(text),
	)
}

func FromFrag(frg frag.Frag) drawable.Unit {
	return Wrap(
		drawable_drain.UnitFromFrags(frg),
	)
}

func (u *InputLineUnit) ToUnit() drawable.Unit {
	return drawable.NewBuilder().
		Name(Name).
		Boot(u.boot).
		Wipe(u.unit.Drawable.Wipe).
		Draw(u.draw).
		ToUnit()
}

func (u *InputLineUnit) boot() {
	u.loaded = true

	u.unit.Drawable.Boot()
}

func (u *InputLineUnit) draw(size winsize.Winsize) ([]line.Line, bool) {
	assert.True(u.loaded, drawable.MessageInitialized)

	if size.Rows == 0 {
		return make([]line.Line, 0), false
	}

	lines, _ := drain.UnitLazy(size, u.unit)
	if len(lines) == 0 {
		prompt := line.FromString(u.prompt)
		return []line.Line{prompt}, false
	}

	prompt := frag.FromString(
		u.prompt + marker.DefaultPaddingText,
	)

	lines[0] = line.BuilderFromLine(lines[0]).
		UnshiftFrags(prompt).
		Line()

	return lines, false
}
