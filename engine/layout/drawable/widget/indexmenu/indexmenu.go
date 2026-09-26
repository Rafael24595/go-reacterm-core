package indexmenu

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/format"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/drain"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

const Name = "index_menu_unit"

type IndexMenuUnit struct {
	loaded  bool
	pointer Pointer
	indexer format.IndexProvider
	options []frag.Frag
	cursor  uint16
	unit    drawable.Unit
}

func New(options []frag.Frag) *IndexMenuUnit {
	clone := make([]frag.Frag, len(options))
	copy(clone, options)

	return &IndexMenuUnit{
		loaded:  false,
		pointer: pointerSelect,
		indexer: format.HyphenIndex,
		options: clone,
		cursor:  0,
		unit:    drawable.Unit{},
	}
}

func UnitFromOptions(options []frag.Frag) drawable.Unit {
	return New(options).ToUnit()
}

func (u *IndexMenuUnit) Pointer(pointer Pointer) *IndexMenuUnit {
	u.pointer = pointer
	return u
}

func (u *IndexMenuUnit) Indexer(indexer format.IndexProvider) *IndexMenuUnit {
	u.indexer = indexer
	return u
}

func (u *IndexMenuUnit) Cursor(cursor uint16) *IndexMenuUnit {
	u.cursor = cursor
	return u
}

func (u *IndexMenuUnit) ToUnit() drawable.Unit {
	return drawable.NewBuilder().
		Name(Name).
		Boot(u.boot).
		Wipe(u.wipe).
		Draw(u.draw).
		ToUnit()
}

func (u *IndexMenuUnit) boot() {
	u.loaded = true

	lines := make([]line.Line, 0)

	digits := math.Digits(len(u.options))
	generator := u.indexer(winsize.Cols(digits))

	for i, o := range u.options {
		focusAtom := atom.None
		selectAtom := atom.None
		if i == int(u.cursor) {
			focusAtom = atom.Focus
			if u.pointer == pointerSelect {
				selectAtom = atom.Select
			}
		}

		alignFrag := frag.FromSpec(spec.JustifyRight(2))

		indexFrag := u.makeIndex(i, generator).
			AddAtom(selectAtom).
			Frag()

		spacerFrag := frag.TextAtom(
			marker.DefaultPaddingText, selectAtom,
		)

		titleFrag := frag.TextAtom(
			o.Text(), focusAtom, selectAtom,
		)

		lines = append(lines,
			line.FromFrags(
				alignFrag,
				indexFrag,
				spacerFrag,
				titleFrag,
			),
		)
	}

	unit := drain.UnitFromLines(lines...)
	unit.Drawable.Boot()

	u.unit = unit
}

func (u *IndexMenuUnit) makeIndex(cursor int, generator format.Index) *frag.Builder {
	txt := generator(int(u.cursor), cursor)

	index := frag.NewBuilder(txt)
	if u.pointer == pointerBold && cursor == int(u.cursor) {
		index.AddAtom(atom.Bold)
	}

	return index
}

func (u *IndexMenuUnit) wipe() {
	if u.unit.Drawable.Wipe == nil {
		return
	}
	u.unit.Drawable.Wipe()
}

func (u *IndexMenuUnit) draw(size winsize.Winsize) ([]line.Line, bool) {
	assert.True(u.loaded, drawable.MessageInitialized)

	return u.unit.Drawable.Draw(size)
}
