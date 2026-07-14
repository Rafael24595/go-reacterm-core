package selection

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/model/ascii"
	"github.com/Rafael24595/go-reacterm-core/engine/model/input"
	"github.com/Rafael24595/go-reacterm-core/engine/model/offset"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
)

type Result struct {
	Frags []frag.Frag
	End   offset.Offset
}

type Renderer struct {
	buffer []rune
	start  offset.Offset
	end    offset.Offset
	blink  atom.Atom
}

func NewRenderer(
	buffer []rune,
	start, end offset.Offset,
	blink ...atom.Atom,
) Renderer {
	return Renderer{
		buffer: buffer,
		start:  start,
		end:    end,
		blink:  atom.Merge(blink...),
	}
}

func (r Renderer) selection() []rune {
	return r.buffer[r.start:r.end]
}

func (r Renderer) Resolve(caret *input.TextCursor) Result {
	selection := r.selection()
	if len(selection) == 0 {
		return r.resolveEmpty()
	}

	if caret.Caret() != caret.Anchor() && r.end == caret.Anchor() {
		return r.resolveBackward()
	}

	return r.resolveForward()
}

func (r Renderer) resolveBackward() Result {
	frags := make([]frag.Frag, 0, 2)
	focusAtom := atom.Focus

	selection := r.selection()
	if r.start > 0 && selection[0] == ascii.ENTER_LF {
		focusAtom = atom.None

		frags = append(frags,
			*frag.FromRunes(marker.PrintableCaretRunes).
				AddAtom(r.blink, atom.Focus),
		)
	}

	frags = append(frags,
		*frag.FromRunes(selection).
			AddAtom(r.blink, focusAtom),
	)

	return Result{
		Frags: frags,
		End:   r.end,
	}
}

func (r Renderer) resolveForward() Result {
	selection := r.selection()
	if selection[len(selection)-1] == ascii.ENTER_LF {
		return r.resolveForwardEnter()
	}

	return r.resolveForwardNonEnter()
}

func (r Renderer) resolveForwardNonEnter() Result {
	frags := make([]frag.Frag, 0, 3)

	selection := r.selection()
	if len(selection) > 1 {
		frags = append(frags,
			*frag.FromRunes(selection[:len(selection)-1]).
				AddAtom(r.blink),
		)
	}

	frags = append(frags,
		*frag.FromRunes(selection[len(selection)-1:]).
			AddAtom(r.blink, atom.Focus),
	)

	return Result{
		Frags: frags,
		End:   r.end,
	}
}

func (r Renderer) resolveForwardEnter() Result {
	frags := make([]frag.Frag, 0, 3)

	selection := r.selection()
	if len(selection) == 1 {
		frags = append(frags,
			*frag.FromRunes(marker.PrintableCaretRunes).
				AddAtom(r.blink),
		)
	}

	footer, nextEnd := r.resolveEnterFooter()

	frags = append(frags,
		*frag.FromRunes(selection).
			AddAtom(r.blink),
		*frag.FromRunes(footer).
			AddAtom(r.blink, atom.Focus),
	)

	return Result{
		Frags: frags,
		End:   nextEnd,
	}
}

func (r Renderer) resolveEnterFooter() ([]rune, offset.Offset) {
	if int(r.end) >= len(r.buffer) {
		return marker.PrintableCaretRunes, r.end
	}

	if r.buffer[r.end] == ascii.ENTER_LF {
		return marker.PrintableCaretRunes, r.end
	}

	return r.buffer[r.end : r.end+1], r.end + 1
}

func (r Renderer) resolveEmpty() Result {
	assert.Unreachable("selection should have at least one character")

	frags := []frag.Frag{
		*frag.Empty().AddAtom(atom.Focus),
	}

	return Result{
		Frags: frags,
		End:   r.end,
	}
}
