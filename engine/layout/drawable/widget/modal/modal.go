package modal

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/config/padding/cols"
	"github.com/Rafael24595/go-reacterm-core/engine/config/padding/rows"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/decorator/box"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/spatial/justify"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/spatial/stack"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/drain"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/padding"
	"github.com/Rafael24595/go-reacterm-core/engine/model/hint"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

const Name = "modal_unit"

type ModalUnit struct {
	loaded     bool
	lazyLoaded bool
	text       []line.Line
	options    []frag.Frag
	limit      uint
	cursor     uint16
	unit       drawable.Unit
}

func New() *ModalUnit {
	return &ModalUnit{
		loaded:     false,
		lazyLoaded: false,
		text:       make([]line.Line, 0),
		options:    make([]frag.Frag, 0),
		limit:      style.DefaultMaxOpts,
		cursor:     0,
		unit:       drawable.Unit{},
	}
}

func (u *ModalUnit) AddText(text ...line.Line) *ModalUnit {
	u.text = append(u.text, text...)
	return u
}

func (u *ModalUnit) AddOptions(options ...frag.Frag) *ModalUnit {
	u.options = append(u.options, options...)
	return u
}

func (u *ModalUnit) SetLimit(limit uint) *ModalUnit {
	u.limit = limit
	return u
}

func (u *ModalUnit) SetCursor(cursor uint16) *ModalUnit {
	u.cursor = cursor
	return u
}

func (u *ModalUnit) ToUnit() drawable.Unit {
	return drawable.NewBuilder().
		Name(Name).
		Boot(u.boot).
		Wipe(u.wipe).
		Draw(u.draw).
		ToUnit()
}

func (u *ModalUnit) boot() {
	u.loaded = true
	u.lazyLoaded = false
}

func (u *ModalUnit) lazyBoot(size winsize.Winsize) {
	if u.lazyLoaded {
		return
	}

	u.lazyLoaded = true

	opts := make([]frag.Frag, len(u.options))
	for i := range u.options {
		old := u.options[i]
		opts[i] = *frag.New(old.Text).
			AddAtom(old.Atom).
			AddSpec(old.Spec)

		if i == int(u.cursor) {
			opts[i].AddAtom(atom.Select)
		}
	}

	measure := line.MaxLineMeasure(size.Cols, u.text...) + 1
	text := formatLines(u.text...)

	title := drain.UnitFromLines(text...)

	optionsBlock := drain.Unit(
		justify.New(opts).
			MaxCols(measure).
			ToUnit(),
	)

	title.Drawable.Boot()
	optionsBlock.Drawable.Boot()

	stack := stack.VStackFromUnits(
		title,
		optionsBlock,
	)

	box := box.New(stack).
		PaddingX(
			hint.Fixed[winsize.Cols](1),
		).
		PaddingY(
			hint.Fixed[winsize.Rows](1),
		).
		ToUnit()

	position := padding.NewBuilder().
		Rows(
			hint.Maximize[winsize.Rows](),
			rows.WithPosition(style.Middle),
		).
		Cols(
			hint.Maximize[winsize.Cols](),
			cols.WithPosition(style.Center),
		).
		ToUnit(box)

	position.Drawable.Boot()

	u.unit = position
}

func (u *ModalUnit) wipe() {
	u.lazyLoaded = false

	if u.unit.Drawable.Wipe == nil {
		return
	}

	u.unit.Drawable.Wipe()
}

func (u *ModalUnit) draw(size winsize.Winsize) ([]line.Line, bool) {
	assert.True(u.loaded, drawable.MessageInitialized)

	u.lazyBoot(size)

	return u.unit.Drawable.Draw(size)
}

func formatLines(lines ...line.Line) []line.Line {
	out := make([]line.Line, len(lines))
	copy(out, lines)

	out = append(out, *line.Empty())

	return out
}
