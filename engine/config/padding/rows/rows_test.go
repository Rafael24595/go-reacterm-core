package rows

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
)

func TestResolveConfigDefaults(t *testing.T) {
	cfg := ResolveConfig()
	assert.Equal(t, style.Top, cfg.Position)

	frag := cfg.Provider(
		winsize.New(10, 20),
	)

	assert.Size(t, 0, frag.Text)
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

	lines := []text.Line{
		*text.NewLine("Golang"),
	}

	frag := cfg.Provider(
		winsize.New(10, 20),
		lines...,
	)

	assert.Equal(t, ".", frag.Text)
	assert.True(t, frag.Spec.Kind().HasAny(spec.KindExtendRight))
	assert.Equal(t, "6", frag.Spec.Args()[spec.KeyExtendRightSize].Text())
}
