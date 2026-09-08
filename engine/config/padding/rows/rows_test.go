package rows

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
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

	assert.Empty(t, frag.Text())
}

func TestResolveConfigWithMultipleOptions(t *testing.T) {
	cfg := ResolveConfig(
		WithPosition(style.Bottom),
		WithFrag(frag.FromString("custom")),
	)

	assert.Equal(t, style.Bottom, cfg.Position)

	f := cfg.Provider(winsize.New(10, 20))
	
	assert.Equal(t, "custom", f.Text())
}

func TestWithPosition(t *testing.T) {
	cfg := defaultConfig()

	WithPosition(style.Bottom)(&cfg)

	assert.Equal(t, style.Bottom, cfg.Position)
}

func TestWithFrag(t *testing.T) {
	cfg := defaultConfig()

	WithFrag(
		frag.FromString("golang"),
	)(&cfg)

	frag := cfg.Provider(
		winsize.New(10, 20),
	)

	assert.Equal(t, "golang", frag.Text())
}

func TestWithFillFrag(t *testing.T) {
	t.Run("applies specified fill text", func(t *testing.T) {
		cfg := defaultConfig()
		WithFillFrag(".")(&cfg)

		lines := []line.Line{
			line.FromString("Golang"),
		}

		f := cfg.Provider(winsize.New(10, 20), lines...)

		assert.Equal(t, ".", f.Text())
		assert.True(t, f.Spec().Kind().HasAny(spec.KindExtendRight))
		assert.Equal(t, "6", f.Spec().Args()[spec.KeyExtendRightSize].Text())
	})

	t.Run("uses default padding when empty parameter provided", func(t *testing.T) {
		cfg := defaultConfig()
		WithFillFrag("")(&cfg)

		lines := []line.Line{
			line.FromString("Golang"),
		}

		f := cfg.Provider(winsize.New(10, 20), lines...)

		assert.Equal(t, marker.DefaultPaddingText, f.Text())
	})
}
