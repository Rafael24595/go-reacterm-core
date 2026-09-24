package transformer

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
)

func TestWithMargin_ShouldSubtractDimensions(t *testing.T) {
	initial := winsize.New(30, 100)
	tf := WithMargin(5, 10)

	result := tf(initial)

	assert.Equal(t, 25, result.Rows)
	assert.Equal(t, 90, result.Cols)
}

func TestWithMargin_ShouldClampToZero(t *testing.T) {
	initial := winsize.New(5, 10)
	tf := WithMargin(10, 20)

	result := tf(initial)

	assert.Equal(t, 0, result.Rows)
	assert.Equal(t, 0, result.Cols)
}

func TestCompose_ShouldApplyInOrder(t *testing.T) {
	initial := winsize.New(50, 100)

	composed := Compose(
		WithMargin(5, 10),
		WithMargin(10, 20),
	)

	result := composed(initial)

	assert.Equal(t, 35, result.Rows)
	assert.Equal(t, 70, result.Cols)
}
