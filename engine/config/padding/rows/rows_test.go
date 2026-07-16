package rows

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

func TestResolveConfigDefaults(t *testing.T) {
	cfg := ResolveConfig()
	assert.Equal(t, style.Top, cfg.Position)

	frag := cfg.Provider(
		winsize.New(10, 20),
	)

	assert.Empty(t, frag.Text)
}

func TestWithPosition(t *testing.T) {
	cfg := defaultConfig()

	WithPosition(style.Bottom)(&cfg)

	assert.Equal(t, style.Bottom, cfg.Position)
}

func TestWithFrag(t *testing.T) {
	cfg := defaultConfig()

	WithFrag(
		*frag.New("golang"),
	)(&cfg)

	frag := cfg.Provider(
		winsize.New(10, 20),
	)

	assert.Equal(t, "golang", frag.Text)
}

func TestWithFillFrag(t *testing.T) {
	cfg := defaultConfig()

	WithFillFrag(".")(&cfg)

	lines := []line.Line{
		*line.New("Golang"),
	}

	frag := cfg.Provider(
		winsize.New(10, 20),
		lines...,
	)

	assert.Equal(t, ".", frag.Text)
	assert.True(t, frag.Spec.Kind().HasAny(spec.KindExtendRight))
	assert.Equal(t, "6", frag.Spec.Args()[spec.KeyExtendRightSize].Text())
}
