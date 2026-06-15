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
	options []text.Fragment
	cursor  uint16
	unit    drawable.Unit
}

func New(options []text.Fragment) *IndexMenuUnit {
	clone := make([]text.Fragment, len(options))
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

func UnitFromOptions(options []text.Fragment) drawable.Unit {
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
		Init(u.init).
		Wipe(u.wipe).
		Draw(u.draw).
		ToUnit()
}

func (u *IndexMenuUnit) init() {
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

		alignFrag := text.EmptyFragment().
			AddSpec(spec.JustifyRight(2))

		indexFrag := u.makeIndex(i, winsize.Cols(digits)).
			AddAtom(selectAtom)

		spacerFrag := text.NewFragment(marker.DefaultPaddingText).
			AddAtom(selectAtom)

		titleFrag := text.NewFragment(o.Text).
			AddAtom(focusAtom, selectAtom)

		lines = append(lines,
			*text.LineFromFragments(
				*alignFrag,
				*indexFrag,
				*spacerFrag,
				*titleFrag,
			),
		)
	}

	unit := drain.UnitFromLines(lines...)
	unit.Drawable.Init()

	u.unit = unit
}

func (u *IndexMenuUnit) makeIndex(cursor int, digits winsize.Cols) *text.Fragment {
	if u.meta.Kind == marker.Numeric {
		return u.makeNumericIndex(cursor, digits)
	}

	if u.meta.Kind == marker.Alphabetic {
		return u.makeAlphabeticIndex(cursor, digits)
	}

	return u.makeCustomIndex(cursor)
}

func (u *IndexMenuUnit) makeCustomIndex(cursor int) *text.Fragment {
	data := u.meta.Index
	if cursor == int(u.cursor) {
		data = u.meta.Cursor
	}

	return text.NewFragment(data)
}

func (u *IndexMenuUnit) makeNumericIndex(cursor int, digits winsize.Cols) *text.Fragment {
	text := format.TextFromAny(cursor + 1)
	return u.makeTextIndex(cursor, digits, text)
}

func (u *IndexMenuUnit) makeAlphabeticIndex(cursor int, digits winsize.Cols) *text.Fragment {
	text := format.TextFromAny(
		format.NumberToAlpha(cursor),
	)
	return u.makeTextIndex(cursor, digits, text)
}

func (u *IndexMenuUnit) makeTextIndex(cursor int, digits winsize.Cols, text format.Text) *text.Fragment {
	filler := marker.DefaultPaddingText
	data := format.JustifyLeft(digits, text, filler)
	return u.makeCommonIndex(cursor, data)
}

func (u *IndexMenuUnit) makeCommonIndex(cursor int, txt string) *text.Fragment {
	index := text.NewFragment(txt + ".- ")
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
