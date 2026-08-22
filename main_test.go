package main

import (
	"testing"

	"github.com/Rafael24595/go-reacterm-core/engine/app/hash"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/partial/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/partial/pipeline/inline"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/partial/pipeline/spacer"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/wrapper/cache"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/wrapper/help"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/wrapper/history"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/wrapper/pagination"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/layout"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/composer"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render"
	"github.com/Rafael24595/go-reacterm-core/engine/render/chunk"
	"github.com/Rafael24595/go-reacterm-core/engine/render/processor"
	"github.com/Rafael24595/go-reacterm-core/engine/render/styler"
	"github.com/Rafael24595/go-reacterm-core/engine/render/wrap"
	"github.com/Rafael24595/go-reacterm-core/engine/render/wrap/delta"

	system_cache "github.com/Rafael24595/go-reacterm-core/engine/commons/structure/cache"
	wrap_processor "github.com/Rafael24595/go-reacterm-core/engine/render/wrap/processor"
	wrap_splitter "github.com/Rafael24595/go-reacterm-core/engine/render/wrap/splitter"
	wrapper_render "github.com/Rafael24595/go-reacterm-core/wrapper/render"
	wrapper_screen "github.com/Rafael24595/go-reacterm-core/wrapper/screen"
)

func BenchmarkPipeline_WithOutCache(b *testing.B) {
	size := winsize.New(80, 200)

	node := helperMakeNodeWithOutCache()
	layout := helperMakeLayout()
	render := helperMakeRender()

	ui := state.NewUIState()

	node.Screen.Boot(*ui)

	for b.Loop() {
		ui.Pager.TargetPage += 1

		vm := node.Screen.View(*ui)
		_, lines := layout.Compose(ui, vm, size)
		_ = render.Processor(lines, size)
	}
}

func BenchmarkPipeline_WithCache(b *testing.B) {
	size := winsize.New(80, 200)

	che := system_cache.NewMemory[hash.Hash, delta.Delta]()

	wrap.DefineWrapper(
		helperMakeWrapper(che),
	)

	node := helperMakeNodeWithCache(che)
	layout := helperMakeLayout()
	render := helperMakeRender()

	ui := state.NewUIState()

	node.Screen.Boot(*ui)

	for b.Loop() {
		ui.Pager.TargetPage += 1

		vm := node.Screen.View(*ui)
		_, lines := layout.Compose(ui, vm, size)
		_ = render.Processor(lines, size)
	}
}

func helperMakeNodeWithOutCache() screen.Node {
	landing := wrapper_screen.NewDemoTextArea()

	history := history.New(landing).ToNode()
	pagination := pagination.New(history).
		ToNode()
	helper := help.New(pagination).ToNode()

	return helperMakePipeline(helper)
}

func helperMakeNodeWithCache(systemCache system_cache.Cache[hash.Hash, delta.Delta]) screen.Node {
	landing := wrapper_screen.NewDemoTextArea()

	cache := cache.New(systemCache, landing).ToNode()
	history := history.New(cache).ToNode()
	pagination := pagination.New(history).
		ToNode()
	helper := help.New(pagination).ToNode()

	return helperMakePipeline(helper)
}

func helperMakeLayout() layout.Layout {
	return layout.NewBuilder(composer.Standard).
		ToLayout()
}

func helperMakeRender() render.Render {
	atomStyler := styler.NewDefaultAtom().
		Push(wrapper_render.Atoms...)

	specStyler := styler.NewDefaultSpec()

	standard := processor.New(*atomStyler, *specStyler)

	return render.NewBuilder(standard.Render).
		ToRender()
}

func helperMakeWrapper(cache system_cache.Cache[hash.Hash, delta.Delta]) wrap.Wrapper {
	return wrap.FromWrapper(
		wrap.DefaultWrapper(),
		wrap.WithProcessors(
			wrap_processor.Chunk(chunk.DefaultChunk),
		),
		wrap.WithSplitter(
			wrap_splitter.SplitLineWithCache(cache),
		),
	)
}

func helperMakePipeline(node screen.Node) screen.Node {
	headerStep := wrapper_screen.NewDemoHeader()

	inlineStep := inline.Transformer(
		inline.DefaultSeparator,
		pipeline.NewFilter(pipeline.Tags, screen.SystemMetaTag),
		pipeline.Footer,
		pipeline.After,
	)

	spacerHeader := spacer.Transformer(
		spacer.NewMeta(1, spacer.Between, pipeline.After),
		pipeline.Header,
	)

	spacerFooter := spacer.Transformer(
		spacer.NewMeta(1, spacer.Between, pipeline.Before),
		pipeline.Footer,
	)

	return pipeline.New(node,
		headerStep, inlineStep, spacerHeader, spacerFooter,
	).ToNode()
}
