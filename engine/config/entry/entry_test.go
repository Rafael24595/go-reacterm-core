package entry

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/config/layer"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"

	screen_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/screen"
)

func TestNewEntry(t *testing.T) {
	t.Run("successfully constructs entry", func(t *testing.T) {
		mock := screen_test.MockNode{}
		var dummyOpt layer.Option[winsize.Rows]

		e := New(mock.ToNode(), Selectable(), WithLayout(dummyOpt))

		assert.True(t, e.Selectable)
		assert.Size(t, 1, e.Opts)
	})

	t.Run("panics on zero node", func(t *testing.T) {
		assert.Panic(t, func() {
			New(screen.Node{})
		})
	})
}
