package textarea

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/primitive/lines"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/widget/textarea/selection"
	"github.com/Rafael24595/go-reacterm-core/engine/model/input"
	"github.com/Rafael24595/go-reacterm-core/engine/model/offset"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
	"github.com/Rafael24595/go-reacterm-core/engine/render/wrap"
)

const Name = "text_area_unit"

type TextAreaUnit struct {
	loaded     bool
	lazyLoaded bool
	writeMode  bool
	indexMode  bool
	buffer     []rune
	caret      *input.TextCursor
	steps      []Transformer
	unit       drawable.Unit
}

func New(buffer []rune, caret *input.TextCursor) *TextAreaUnit {
	clone := make([]rune, len(buffer))
	copy(clone, buffer)

	return &TextAreaUnit{
		loaded:     false,
		lazyLoaded: false,
		writeMode:  false,
		indexMode:  false,
		buffer:     clone,
		caret:      caret,
		steps:      make([]Transformer, 0),
		unit:       drawable.Unit{},
	}
}

func (u *TextAreaUnit) WriteMode(writeMode bool) *TextAreaUnit {
	u.writeMode = writeMode
	return u
}

func (u *TextAreaUnit) IndexMode(indexMode bool) *TextAreaUnit {
	u.indexMode = indexMode
	return u
}

func (u *TextAreaUnit) PushStep(step Transformer) *TextAreaUnit {
	u.steps = append(u.steps, step)
	return u
}

func (u *TextAreaUnit) ToUnit() drawable.Unit {
	return drawable.NewBuilder().
		Name(Name).
		Boot(u.boot).
		Wipe(u.wipe).
		Draw(u.draw).
		ToUnit()
}

func (u *TextAreaUnit) boot() {
	u.loaded = true
	u.lazyLoaded = false
}

func (u *TextAreaUnit) lazyBoot(size winsize.Winsize) {
	if u.lazyLoaded {
		return
	}

	u.lazyLoaded = true

	start := u.caret.SelectStart().Sub(1)
	end := u.caret.SelectEnd()

	if len(u.buffer) == 0 {
		u.buffer = append(u.buffer, marker.PrintableCaretRunes...)
		start = 0
		end = 1
	}

	frags := u.resolveFrags(u.buffer, start, end)
	for _, step := range u.steps {
		frags = step(frags)
	}

	base := line.FromFrags(frags...)

	result := u.makeLines(base)
	result = wrap.MaterializeEmpty(size, marker.DefaultPaddingText, result...)

	unit := lines.UnitFromLayout(result...)
	unit.Drawable.Boot()

	u.unit = unit
}

func (u *TextAreaUnit) makeLines(base *line.Line) []wrap.LayoutLine {
	if u.indexMode {
		return wrap.NormalizeLinesWithOrder(*base)
	}
	return wrap.NormalizeLines(*base)

}

func (u *TextAreaUnit) wipe() {
	u.lazyLoaded = false

	if u.unit.Drawable.Wipe == nil {
		return
	}

	u.unit.Drawable.Wipe()
}

func (u *TextAreaUnit) resolveFrags(
	renderBuffer []rune,
	start, end offset.Offset,
) []frag.Frag {
	frags := make([]frag.Frag, 0, 6)

	bufferLen := offset.Offset(len(renderBuffer))

	start = min(start, bufferLen.Sub(1))
	end = min(end, bufferLen)

	if start > 0 {
		frags = append(frags,
			*frag.New(string(renderBuffer[:start])),
		)
	}

	renderer := selection.NewRenderer(
		renderBuffer, start, end, u.blinkStyle(),
	)

	result := renderer.Resolve(u.caret)

	end = result.End
	frags = append(frags, result.Frags...)

	if int(end) < len(renderBuffer) {
		frags = append(frags,
			*frag.New(string(renderBuffer[end:])),
		)
	}

	return frags
}

func (u *TextAreaUnit) blinkStyle() atom.Atom {
	if !u.writeMode {
		return atom.None
	}

	return u.caret.BlinkStyle()
}

func (u *TextAreaUnit) draw(size winsize.Winsize) ([]line.Line, bool) {
	assert.True(u.loaded, drawable.MessageInitialized)

	u.lazyBoot(size)

	return u.unit.Drawable.Draw(size)
}
