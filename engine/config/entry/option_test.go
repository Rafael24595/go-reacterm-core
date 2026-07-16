package entry

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-reacterm-core/engine/config/layer"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	screen_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/screen"
)

func TestDefaultEntry(t *testing.T) {
	mock := screen_test.MockNode{}
	cfg := defaultEntry(mock.ToNode())

	assert.False(t, cfg.Selectable)
	assert.NotNil(t, cfg.Opts)
	assert.Empty(t, cfg.Opts)
}

func TestSelectableOption(t *testing.T) {
	mock := screen_test.MockNode{}
	cfg := defaultEntry(mock.ToNode())

	opt := Selectable()
	opt(&cfg)

	assert.True(t, cfg.Selectable)
}

func TestWithLayoutOption(t *testing.T) {
	mock := screen_test.MockNode{}
	cfg := defaultEntry(mock.ToNode())

	var dummyOpt1 layer.Option[winsize.Rows]
	var dummyOpt2 layer.Option[winsize.Rows]

	opt := WithLayout(dummyOpt1, dummyOpt2)
	opt(&cfg)

	assert.Size(t, 2, cfg.Opts)
}
