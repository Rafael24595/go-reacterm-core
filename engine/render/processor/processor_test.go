package processor

import (
	"strings"
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

func TestWithPadding(t *testing.T) {
	size := winsize.New(5, 10)

	mockInner := func(lines []line.Line, innerSize winsize.Winsize) []string {
		return []string{"OK"}
	}

	t.Run("applies correct padding when reducing inner dimensions", func(t *testing.T) {
		reduceTransform := func(_ winsize.Winsize) winsize.Winsize {
			return winsize.New(1, 4)
		}

		paddedProc := WithPadding(reduceTransform, mockInner)
		result := paddedProc([]line.Line{line.FromString("OK")}, size)

		lines := strings.Split(result, "\n")

		assert.Equal(t, 5, len(lines))
		assert.Equal(t, "          ", lines[0])
		assert.Equal(t, "          ", lines[1])

		assert.Equal(t, "   OK     ", lines[2])

		assert.Equal(t, "          ", lines[3])
		assert.Equal(t, "          ", lines[4])
	})

	t.Run("clamps target dimensions to prevent canvas overflow", func(t *testing.T) {
		overflowTransform := func(_ winsize.Winsize) winsize.Winsize {
			return winsize.New(10, 20)
		}

		paddedProc := WithPadding(overflowTransform, mockInner)
		result := paddedProc([]line.Line{line.FromString("OK")}, size)

		lines := strings.Split(result, "\n")

		assert.Equal(t, 5, len(lines))
		assert.Equal(t, "OK        ", lines[0])
	})
}
