package drawable

import (
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"
)

const (
	ErrorMissingName = "missing_name"
	ErrorMissingBoot = "missing_boot"
	ErrorMissingWipe = "missing_wipe"
	ErrorMissingDraw = "missing_draw"
)

type Builder struct {
	name string
	tags set.Set[string]
	boot BootFunc
	wipe WipeFunc
	draw DrawFunc
}

func NewBuilder() *Builder {
	return &Builder{
		name: "",
		tags: set.New[string](),
		boot: nil,
		wipe: nil,
		draw: nil,
	}
}

func (b *Builder) Name(name string) *Builder {
	b.name = name
	return b
}

func (b *Builder) AddTags(tags ...string) *Builder {
	b.tags.Add(tags...)
	return b
}

func (b *Builder) MergeTags(tags set.Set[string]) *Builder {
	b.tags.Merge(tags)
	return b
}

func (b *Builder) Boot(boot BootFunc) *Builder {
	b.boot = boot
	return b
}

func (b *Builder) Wipe(wipe WipeFunc) *Builder {
	b.wipe = wipe
	return b
}

func (b *Builder) Draw(draw DrawFunc) *Builder {
	b.draw = draw
	return b
}

func (b *Builder) makeTags() set.Set[string] {
	tags := set.New[string]()

	if b.name == "" {
		tags.Add(ErrorMissingName)
	}

	if b.boot == nil {
		b.tags.Add(ErrorMissingBoot)
	}

	if b.wipe == nil {
		b.tags.Add(ErrorMissingWipe)
	}

	if b.draw == nil {
		b.tags.Add(ErrorMissingDraw)
	}

	tags.Merge(b.tags)

	return tags
}

func (b *Builder) toDrawable() Drawable {
	return Drawable{
		Boot: b.boot,
		Wipe: b.wipe,
		Draw: b.draw,
	}
}

func (b *Builder) ToUnit() Unit {
	return Unit{
		Name:     b.name,
		Tags:     b.makeTags(),
		Drawable: b.toDrawable(),
	}
}
