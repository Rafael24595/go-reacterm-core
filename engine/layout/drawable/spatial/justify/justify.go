package justify

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

const (
	SlotsBetween = 0
	SlotsEvenly  = 1
	SlotsAround  = 2
)

const Name = "justify_unit"

type JustifyUnit struct {
	loaded  bool
	maxOpts uint16
	maxCols winsize.Cols
	justify style.Justify
	frags   []frag.Frag
	cursor  uint16
}

func New(frags []frag.Frag) *JustifyUnit {
	return &JustifyUnit{
		loaded:  false,
		maxOpts: style.DefaultMaxOpts,
		justify: style.JustifyAround,
		frags:   frags,
		cursor:  0,
	}
}

func UnitFromFrags(frags []frag.Frag) drawable.Unit {
	return New(frags).ToUnit()
}

func (u *JustifyUnit) MaxOpts(opts uint16) *JustifyUnit {
	u.maxOpts = max(1, opts)
	return u
}

func (u *JustifyUnit) MaxCols(cols winsize.Cols) *JustifyUnit {
	u.maxCols = max(1, cols)
	return u
}

func (u *JustifyUnit) Justify(justify style.Justify) *JustifyUnit {
	u.justify = justify
	return u
}

func (u *JustifyUnit) AddFrags(frags []frag.Frag) *JustifyUnit {
	u.frags = append(u.frags, frags...)
	return u
}

func (u *JustifyUnit) ToUnit() drawable.Unit {
	return drawable.NewBuilder().
		Name(Name).
		Boot(u.boot).
		Wipe(u.wipe).
		Draw(u.draw).
		ToUnit()
}

func (u *JustifyUnit) boot() {
	u.loaded = true

	u.cursor = 0
}

func (u *JustifyUnit) wipe() {
	u.cursor = 0
}

func (u *JustifyUnit) draw(size winsize.Winsize) ([]line.Line, bool) {
	assert.True(u.loaded, drawable.MessageInitialized)

	if u.cursor >= uint16(len(u.frags)) {
		return make([]line.Line, 0), false
	}

	maxOpts := int(u.maxOpts)
	maxCols := math.MinNotZero(size.Cols, u.maxCols)

	remaining := winsize.Cols(0)
	frags := make([]frag.Frag, 0)

	for i := u.cursor; i < uint16(len(u.frags)); i++ {
		frg := u.frags[i]

		fragsLen := len(frags)
		fragSize := frag.Measure(size.Cols, frg)

		spacing := winsize.Cols(0)
		if fragsLen > 0 {
			spacing = 1
		}

		newRemaining := remaining + spacing + fragSize
		if fragsLen > 0 && fragsLen >= maxOpts || newRemaining > maxCols {
			justify := justifyLine(maxCols, frags, remaining, u.justify)
			return []line.Line{*justify}, true
		}

		remaining = newRemaining
		frags = append(frags, frg)

		u.cursor += 1
	}

	justify := justifyLine(maxCols, frags, remaining, u.justify)
	return []line.Line{*justify}, u.cursor < uint16(len(u.frags))
}

func justifyLine(cols winsize.Cols, frags []frag.Frag, size winsize.Cols, mode style.Justify) *line.Line {
	line := line.FromFrags(
		addGaps(cols, frags, size, mode)...,
	)

	switch mode {
	case style.JustifyStart:
		return line.AddSpec(
			spec.AlignLeft(),
		)
	case style.JustifyEnd:
		return line.AddSpec(
			spec.AlignRight(),
		)
	case style.JustifyCenter, style.JustifyAround, style.JustifyEvenly:
		return line.AddSpec(
			spec.AlignCenter(),
		)
	}

	return line
}

func addGaps(
	cols winsize.Cols,
	frags []frag.Frag,
	size winsize.Cols,
	mode style.Justify,
) []frag.Frag {
	if len(frags) == 0 {
		return frags
	}

	out := make([]frag.Frag, len(frags))
	copy(out, frags)

	free := cols.Sub(size)
	gaps := len(out) - 1

	if free <= 0 || gaps <= 0 {
		return addSpaceBetween(out)
	}

	switch mode {

	case style.JustifyBetween:
		return distributeSpace(free, out, SlotsBetween)

	case style.JustifyAround:
		return distributeSpace(free, out, SlotsAround)

	case style.JustifyEvenly:
		return distributeSpace(free, out, SlotsEvenly)
	}

	return addSpaceBetween(out)
}

func distributeSpace(free winsize.Cols, frags []frag.Frag, extraSlots winsize.Cols) []frag.Frag {
	gaps := winsize.Cols(
		max(0, len(frags)-1),
	)

	slots := gaps + extraSlots
	base := free / slots
	remainder := free % slots

	out := make([]frag.Frag, len(frags))
	copy(out, frags)

	fix := winsize.Cols(0)
	for i := range gaps {
		gap := base
		if remainder > 0 {
			gap++
			remainder--
		}

		if gap <= 0 {
			continue
		}

		space := frag.Empty().AddSpec(
			spec.JustifyLeft(gap, marker.DefaultPaddingText),
		)

		at := i + fix + 1

		next := make([]frag.Frag, 0, len(out)+1)

		next = append(next, out[:at]...)
		next = append(next, *space)
		next = append(next, out[at:]...)

		out = next

		fix += 1
	}

	return out
}

func addSpaceBetween(frags []frag.Frag) []frag.Frag {
	spaced := make([]frag.Frag, 0, (len(frags)*2)-1)
	for i, f := range frags {
		spaced = append(spaced, f)
		if i < len(frags)-1 {
			spaced = append(spaced,
				*frag.New(marker.DefaultPaddingText),
			)
		}
	}
	return spaced
}
