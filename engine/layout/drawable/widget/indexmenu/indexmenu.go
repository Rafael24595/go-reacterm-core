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
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

const Name = "index_menu_unit"

type IndexMenuUnit struct {
	loaded  bool
	pointer Pointer
	meta    marker.IndexMeta
	options []text.Frag
	cursor  uint16
	unit    drawable.Unit
}

func New(options []text.Frag) *IndexMenuUnit {
	clone := make([]text.Frag, len(options))
	copy(clone, options)

	return &IndexMenuUnit{
		loaded:  false,
		pointer: pointerSelect,
		meta:    marker.HyphenIndex,
		options: clone,
		cursor:  0,
		unit:    drawable.Unit{},
	}
}

func UnitFromOptions(options []text.Frag) drawable.Unit {
	return New(options).ToUnit()
}

func (u *IndexMenuUnit) Pointer(pointer Pointer) *IndexMenuUnit {
	u.pointer = pointer
	return u
}

func (u *IndexMenuUnit) Meta(meta marker.IndexMeta) *IndexMenuUnit {
	u.meta = meta
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

	lines := make([]text.Line, 0)

	digits := math.Digits(len(u.options))

	for i, o := range u.options {
		focusAtom := atom.None
		selectAtom := atom.None
		if i == int(u.cursor) {
			focusAtom = atom.Focus
			if u.pointer == pointerSelect {
				selectAtom = atom.Select
			}
		}

		alignFrag := text.EmptyFrag().
			AddSpec(spec.JustifyRight(2))

		indexFrag := u.makeIndex(i, winsize.Cols(digits)).
			AddAtom(selectAtom)

		spacerFrag := text.NewFrag(marker.DefaultPaddingText).
			AddAtom(selectAtom)

		titleFrag := text.NewFrag(o.Text).
			AddAtom(focusAtom, selectAtom)

		lines = append(lines,
			*text.LineFromFrags(
				*alignFrag,
				*indexFrag,
				*spacerFrag,
				*titleFrag,
			),
		)
	}

	unit := drain.UnitFromLines(lines...)
	unit.Drawable.Boot()

	u.unit = unit
}

func (u *IndexMenuUnit) makeIndex(cursor int, digits winsize.Cols) *text.Frag {
	if u.meta.Kind == marker.Numeric {
		return u.makeNumericIndex(cursor, digits)
	}

	if u.meta.Kind == marker.Alphabetic {
		return u.makeAlphabeticIndex(cursor, digits)
	}

	return u.makeCustomIndex(cursor)
}

func (u *IndexMenuUnit) makeCustomIndex(cursor int) *text.Frag {
	data := u.meta.Index
	if cursor == int(u.cursor) {
		data = u.meta.Cursor
	}

	return text.NewFrag(data)
}

func (u *IndexMenuUnit) makeNumericIndex(cursor int, digits winsize.Cols) *text.Frag {
	text := format.TextFromAny(cursor + 1)
	return u.makeTextIndex(cursor, digits, text)
}

func (u *IndexMenuUnit) makeAlphabeticIndex(cursor int, digits winsize.Cols) *text.Frag {
	text := format.TextFromAny(
		format.NumberToAlpha(cursor),
	)
	return u.makeTextIndex(cursor, digits, text)
}

func (u *IndexMenuUnit) makeTextIndex(cursor int, digits winsize.Cols, text format.Text) *text.Frag {
	filler := marker.DefaultPaddingText
	data := format.JustifyLeft(digits, text, filler)
	return u.makeCommonIndex(cursor, data)
}

func (u *IndexMenuUnit) makeCommonIndex(cursor int, txt string) *text.Frag {
	index := text.NewFrag(txt + ".- ")
	if u.pointer == pointerBold && cursor == int(u.cursor) {
		index.Atom |= atom.Bold
	}
	return index
}

func (u *IndexMenuUnit) wipe() {
	if u.unit.Drawable.Wipe == nil {
		return
	}
	u.unit.Drawable.Wipe()
}

func (u *IndexMenuUnit) draw(size winsize.Winsize) ([]text.Line, bool) {
	assert.True(u.loaded, drawable.MessageInitialized)

	return u.unit.Drawable.Draw(size)
}
