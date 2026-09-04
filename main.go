package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-log/log"
	"github.com/Rafael24595/go-log/log/provider/file"
	"github.com/Rafael24595/go-log/log/record"
	"github.com/Rafael24595/go-reacterm-core/engine/app/cleaner/composite"
	"github.com/Rafael24595/go-reacterm-core/engine/app/cleaner/store"
	"github.com/Rafael24595/go-reacterm-core/engine/app/core"
	"github.com/Rafael24595/go-reacterm-core/engine/app/hash"
	"github.com/Rafael24595/go-reacterm-core/engine/app/runtime"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/partial/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/partial/pipeline/inline"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/partial/pipeline/spacer"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/wrapper/cache"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/wrapper/help"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/wrapper/history"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/wrapper/pagination"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/pass"
	"github.com/Rafael24595/go-reacterm-core/engine/layout"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/composer"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize/transformer"
	"github.com/Rafael24595/go-reacterm-core/engine/render"
	"github.com/Rafael24595/go-reacterm-core/engine/render/chunk"
	"github.com/Rafael24595/go-reacterm-core/engine/render/processor"
	"github.com/Rafael24595/go-reacterm-core/engine/render/styler"
	"github.com/Rafael24595/go-reacterm-core/engine/render/wrap"
	"github.com/Rafael24595/go-reacterm-core/engine/render/wrap/delta"
	"github.com/Rafael24595/go-reacterm-core/engine/terminal"

	local "github.com/Rafael24595/go-reacterm-core/engine/commons/log"

	system_cache "github.com/Rafael24595/go-reacterm-core/engine/commons/structure/cache"

	wrap_processor "github.com/Rafael24595/go-reacterm-core/engine/render/wrap/processor"
	wrap_splitter "github.com/Rafael24595/go-reacterm-core/engine/render/wrap/splitter"

	wrapper_render "github.com/Rafael24595/go-reacterm-core/wrapper/render"
	wrapper_console "github.com/Rafael24595/go-reacterm-core/wrapper/terminal/console"

	wrapper_screen "github.com/Rafael24595/go-reacterm-core/wrapper/screen"
)

const paddingCols = 10
const paddingRows = 5

// Move main and wrapper packages to a new project
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	configLog(ctx)
	defer local.Writer(os.Stderr, log.OnClose)

	terminal := makeTerminal(ctx)

	transformer := transformer.WithMargin(paddingRows, paddingCols)
	layout := makeLayout(transformer)
	render := makeRender(transformer)

	cleaner := composite.Cleaner(
		store.Cleaner,
	)

	passes := []screen.Pass{
		pass.ValidateStructure,
	}

	cache := system_cache.NewMemory[hash.Hash, delta.Delta]()

	wrap.DefineWrapper(
		makeWrapper(cache),
	)

	screen := makeNode(cache)

	<-core.NewEngine(
		terminal,
		layout,
		render,
		cleaner,
		screen,
	).AddPass(passes...).
		RunWithContext(ctx)
}

func configLog(ctx context.Context) {
	provider := file.FileProvider{
		Session: runtime.Instance.SessionID(),
	}

	if err := log.DefaultFromProvider(ctx, provider); err != nil {
		panic(err.Error())
	}

	assert.DefaultWriter(
		log.WriterFromCategory(record.WARNING),
	)
}

func makeTerminal(ctx context.Context) terminal.Terminal {
	return wrapper_console.NewBuilder().
		Context(ctx).
		Reactive(wrapper_console.DefaultReactiveDuration).
		Color("\x1b[0;32m").
		ToTerminal()
}

func makeLayout(transformer winsize.Transformer) layout.Layout {
	return layout.NewBuilder(composer.Standard).
		Transformer(transformer).
		ToLayout()
}

func makeRender(transformer winsize.Transformer) render.Render {
	atomStyler := styler.NewDefaultAtom().
		Push(wrapper_render.Atoms...)

	specStyler := styler.NewDefaultSpec()

	standard := processor.New(*atomStyler, *specStyler)

	adapter := processor.WithPadding(
		transformer,
		standard.RawRender,
	)

	return render.NewBuilder(adapter).
		ToRender()
}

func makeWrapper(cache system_cache.Cache[hash.Hash, delta.Delta]) wrap.Wrapper {
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

func makeNode(systemCache system_cache.Cache[hash.Hash, delta.Delta]) screen.Node {
	landing := wrapper_screen.NewDemoClip()

	cache := cache.New(systemCache, landing).ToNode()
	history := history.New(cache).ToNode()
	pagination := pagination.New(history).
		ToNode()
	helper := help.New(pagination).ToNode()

	return makePipeline(helper)
}

func makePipeline(node screen.Node) screen.Node {
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
