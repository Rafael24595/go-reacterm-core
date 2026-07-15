package box

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/config/padding/cols"
	"github.com/Rafael24595/go-reacterm-core/engine/config/padding/rows"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/runes"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/margin"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/transform/drain"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/transform/padding"
	"github.com/Rafael24595/go-reacterm-core/engine/model/hint"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
	"github.com/Rafael24595/go-reacterm-core/engine/render/wrap"
)

const Name = "box_unit"

type BoxUnit struct {
	loaded    bool
	paddingY  *hint.Size[winsize.Rows]
	paddingX  *hint.Size[winsize.Cols]
	separator marker.BoxSeparatorMeta
	unit      drawable.Unit
}

func New(unit drawable.Unit) *BoxUnit {
	return &BoxUnit{
		loaded:    false,
		paddingY:  nil,
		paddingX:  nil,
		separator: marker.DefaultBoxSeparator,
		unit:      unit,
	}
}

func Wrap(unit drawable.Unit) drawable.Unit {
	return New(unit).ToUnit()
}

func (u *BoxUnit) Separator(separator marker.BoxSeparatorMeta) *BoxUnit {
	u.separator = separator
	return u
}

func (u *BoxUnit) PaddingY(hint hint.Size[winsize.Rows]) *BoxUnit {
	u.paddingY = &hint
	return u
}

func (u *BoxUnit) PaddingX(hint hint.Size[winsize.Cols]) *BoxUnit {
	u.paddingX = &hint
	return u
}

func (u *BoxUnit) ToUnit() drawable.Unit {
	return drawable.NewBuilder().
		Name(Name).
		MergeTags(u.unit.Tags).
		Boot(u.boot).
		Wipe(u.unit.Drawable.Wipe).
		Draw(u.draw).
		ToUnit()
}

func (u *BoxUnit) boot() {
	u.loaded = true

	u.unit = u.makeUnit()

	u.unit.Drawable.Boot()
}

func (u *BoxUnit) makeUnit() drawable.Unit {
	if u.paddingY == nil && u.paddingX == nil {
		return u.unit
	}

	return margin.NewBuilder().
		Rows(
			u.paddingY,
			rows.WithPosition(style.Middle),
		).
		Cols(
			u.paddingX,
			cols.WithPosition(style.Center),
		).
		ToUnit(u.unit)
}

func (u *BoxUnit) draw(size winsize.Winsize) ([]line.Line, bool) {
	assert.True(u.loaded, drawable.MessageInitialized)

	innerSize := u.computeInnerSize(size)
	lines, hasNext := drain.UnitLazy(innerSize, u.unit)

	styled := u.styleLines(size, lines...)

	return styled, hasNext
}

// TODO: investigate spec overflow.
func (u *BoxUnit) styleLines(size winsize.Winsize, lines ...line.Line) []line.Line {
	vertical := horizontalStaticSize(u.separator)

	maxLine := line.MaxLineMeasure(size.Cols, lines...)
	measure := min(maxLine+vertical, size.Cols)

	specCover := spec.ExtendLeft(measure)
	cover := line.FromFrags(
		*frag.New(u.separator.Top).AddSpec(specCover),
	)

	result := make([]line.Line, 0)

	result = append(result, *cover)

	available := size.Cols.Sub(vertical)

	transformer := padding.Cols(
		hint.Fixed(maxLine),
		cols.WithPosition(style.Center),
		cols.WithText(u.separator.Space),
	)

	for _, lin := range transformer(size, lines) {
		for _, v := range wrap.Line(available, &lin) {
			line := u.wrapLine(v)
			result = append(result, line)
		}
	}

	result = append(result, *cover)

	return result
}

func (u *BoxUnit) wrapLine(line line.Line) line.Line {
	frags := make([]frag.Frag, 0)

	frags = append(frags, *frag.New(u.separator.Left))
	frags = append(frags, line.Text...)
	frags = append(frags, *frag.New(u.separator.Right))

	line.Text = frags

	return line
}

func (u *BoxUnit) computeInnerSize(size winsize.Winsize) winsize.Winsize {
	vertical := winsize.Rows(2)
	rows := size.Rows.Sub(vertical)

	horizontal := horizontalStaticSize(u.separator)
	cols := size.Cols.Sub(horizontal)

	return winsize.New(rows, cols)
}

func horizontalSeparatorSize(separator marker.BoxSeparatorMeta) (winsize.Cols, winsize.Cols) {
	return runes.Measure(separator.Left), runes.Measure(separator.Right)
}

func horizontalStaticSize(separator marker.BoxSeparatorMeta) winsize.Cols {
	left, right := horizontalSeparatorSize(separator)
	return left + right
}
