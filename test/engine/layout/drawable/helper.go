package drawable_test

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

const NameMockUnit = "mock_unit"

type MockUnit struct {
	Name string
	Tags set.Set[string]

	BootCalled uint
	Boot       drawable.BootFunc
	Wipe       drawable.WipeFunc
	WipeCalled uint
	Draw       drawable.DrawFunc
	DrawCalls  uint

	Lines  []text.Line
	queue  []text.Line
	Batch  uint
	Status bool
}

func (m *MockUnit) ToUnit() drawable.Unit {
	name := NameMockUnit
	if m.Name != "" {
		name = m.Name
	}

	return drawable.NewBuilder().
		Name(name).
		MergeTags(m.Tags).
		Boot(
			func() {
				m.BootCalled += 1

				if m.Boot != nil {
					m.Boot()
				}
			},
		).
		Wipe(
			func() {
				m.WipeCalled += 1

				if m.Wipe != nil {
					m.Wipe()
				}
			},
		).
		Draw(
			func(size winsize.Winsize) ([]text.Line, bool) {
				m.DrawCalls += 1

				if m.Draw != nil {
					return m.Draw(size)
				}

				if m.Batch == 0 {
					return m.Lines, m.Status
				}

				if len(m.queue) == 0 {
					m.queue = m.Lines
				}

				limit := min(int(m.Batch), len(m.queue))

				data := m.queue[:limit]
				m.queue = m.queue[limit:]

				return data, len(m.queue) > 0
			},
		).
		ToUnit()
}

func Test_UnitBasicSuite(t *testing.T, unit drawable.Unit) {
	t.Helper()

	Helper_ToUnit(t, unit)
	assert.Panic(t, func() {
		unit.Drawable.Draw(winsize.Winsize{})
	})
}

func Helper_ToUnit(t *testing.T, unit drawable.Unit) {
	t.Helper()

	assert.NotEqual(t, "", unit.Name, "Unit.Name should be set")
	assert.True(t, len(unit.Tags) >= 0, "Unit.Tags should be set")

	assert.NotNil(t, unit.Drawable.Boot, "Drawable.Boot should be set")
	assert.NotNil(t, unit.Drawable.Wipe, "Drawable.Wipe should be set")
	assert.NotNil(t, unit.Drawable, "Drawable.Draw should be set")
}
