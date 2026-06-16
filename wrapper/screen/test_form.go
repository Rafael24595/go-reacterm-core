package wrapper_screen

import (
	"time"

	node_pipeline "github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/partial/pipeline"
	text_screen "github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/primitive/text"
	drawable_pipeline "github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline"

	"github.com/Rafael24595/go-reacterm-core/engine/app/pager/action"
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager/predicate"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/behavior"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/behavior/tick"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/behavior/view"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/partial/form"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/partial/pipeline/page"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/primitive/talk"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/config/entry"
	"github.com/Rafael24595/go-reacterm-core/engine/config/layer"
	"github.com/Rafael24595/go-reacterm-core/engine/config/padding/cols"
	"github.com/Rafael24595/go-reacterm-core/engine/config/padding/rows"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/decorator/box"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/spatial/stack"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/focus"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/padding"
	"github.com/Rafael24595/go-reacterm-core/engine/model/buffer"
	"github.com/Rafael24595/go-reacterm-core/engine/model/buffer/processor"
	"github.com/Rafael24595/go-reacterm-core/engine/model/buffer/rule"
	"github.com/Rafael24595/go-reacterm-core/engine/model/chat"
	"github.com/Rafael24595/go-reacterm-core/engine/model/hint"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

type mockTalkService struct {
	owner string
	text  []rune
	chat  []chat.Message
}

func NewTestForm() screen.Node {
	service := initService()

	node := form.New().
		AddNode(
			makeTalk(service),
			entry.Selectable(),
			entry.WithLayout(
				layer.Percent[winsize.Rows](75),
			),
		).
		AddBreak(1).
		AddNode(
			makeTextArea(service),
			entry.Selectable(),
			entry.WithLayout(
				layer.Static[winsize.Rows](),
			),
		).
		PushSteps(
			page.Use(action.Scroll()),
		).
		ToNode()

	return node
}

func initService() *mockTalkService {
	return &mockTalkService{
		owner: "human_001",
		text:  []rune("Hello Golang!"),
		chat: []chat.Message{
			{
				Time:    time.Now().Add(-15 * time.Minute).UnixMilli(),
				Owner:   "human_001",
				Message: "Lorem ipsum dolor sit amet, consectetur adipiscing elit?",
			},
			{
				Time:    time.Now().Add(-14 * time.Minute).UnixMilli(),
				Owner:   "human_002",
				Message: "Morbi ac ex sit amet diam euismod vulputate ut eu leo.",
			},
			{
				Time:    time.Now().Add(-12 * time.Minute).UnixMilli(),
				Owner:   "human_001",
				Message: "Ok.",
			},
			{
				Time:    time.Now().Add(-10 * time.Minute).UnixMilli(),
				Owner:   "human_001",
				Message: "Nullam quis ante sodales, aliquet turpis ut, suscipit erat. Cras nec viverra dolor, non egestas erat. Vivamus ac pretium lectus. Proin id ligula scelerisque, condimentum elit sit amet, imperdiet magna.",
			},
			{
				Time:    time.Now().Add(-8 * time.Minute).UnixMilli(),
				Owner:   "human_002",
				Message: "Wow! Nunc imperdiet, turpis vel dictum pretium, sem nibh sodales est, nec pulvinar diam leo ac augue.",
			},
			{
				Time:    time.Now().Add(-7 * time.Minute).UnixMilli(),
				Owner:   "human_002",
				Message: "Quisque facilisis nisl nec ex feugiat, non tristique sem finibus.",
			},
			{
				Time:    time.Now().Add(-5 * time.Minute).UnixMilli(),
				Owner:   "human_001",
				Message: "Sed hendrerit elementum lorem, vel interdum velit. Vestibulum rhoncus rhoncus mi, in efficitur elit. Duis imperdiet dictum erat, vel laoreet lorem hendrerit eu.",
			},
			{
				Time:    time.Now().Add(-4 * time.Minute).UnixMilli(),
				Owner:   "human_001",
				Message: "Aenean lacinia porta dictum. Ut sed pulvinar purus, eget pretium tellus. In pretium finibus eros id pretium. Aliquam id interdum magna. Proin feugiat, turpis quis tincidunt elementum, neque justo efficitur elit, ac egestas ex lacus ac ante.",
			},
			{
				Time:    time.Now().Add(-2 * time.Minute).UnixMilli(),
				Owner:   "human_002",
				Message: "Entendido. Proin sollicitudin mi ac arcu dictum, eleifend varius tellus ultrices.",
			},
			{
				Time:    time.Now().Add(-1 * time.Minute).UnixMilli(),
				Owner:   "human_001",
				Message: "Donec id elit non mi porta gravida at eget metus. Nulli magna feugiat purus, ac porttitor elit sem id tellus. Aliquam erat volutpat.",
			},
		},
	}
}

func makeTalk(service *mockTalkService) screen.Node {
	talk := talk.New().
		SetName("talk-form - amet").
		SetOwner("human_001").
		ToNode()

	talk = view.Use(
		talk,
		onTalkView(service),
	)

	return talk
}

func onTalkView(service *mockTalkService) view.Middleware {
	return func(uiState state.UIState, context behavior.Context[screen.ViewFunc]) viewmodel.ViewModel {
		messagesLen := 0
		if messages, ok := talk.KeyMessages.Get(
			uiState.Store,
			context.Target.Name,
		); ok {
			messagesLen = len(messages)
		}

		talk.KeyMessages.Set(
			uiState.Store,
			context.Target.Name,
			service.chat,
		)

		vm := context.Next(uiState)

		if messagesLen != len(service.chat) {
			vm.Pager.SetPredicate(predicate.Last())
		}

		return vm
	}
}

func makeTextArea(service *mockTalkService) screen.Node {
	textscreen := text_screen.NewArea().
		SetName("textarea-form - amet").
		SetBuffer(buffer.NewRuneBuffer().
			PushRules(rule.Full...).
			Processor(processor.Identity)).
		EnableBlinking().
		ToNode()

	textscreen = view.Use(
		textscreen,
		onTextAreaView(service),
	)

	textscreen = tick.Use(
		textscreen,
		onTextAreaTick(service),
	)

	textscreen = tick.OnKey(
		textscreen,
		onKeyEnter(service),
		key.ActionEnter,
	)

	pipeline := node_pipeline.New(textscreen).
		PushSteps(wrapStep).
		ExpireOnNode().
		ToNode()

	return pipeline
}

func onTextAreaTick(service *mockTalkService) tick.Middleware {
	return func(uiState *state.UIState, event screen.Event, context behavior.Context[screen.TickFunc]) screen.Result {
		result := context.Next(uiState, event)

		if state, ok := text_screen.KeyState.Get(
			uiState.Store,
			context.Target.Name,
		); ok {
			service.text = state.Buffer
		}

		return result
	}
}

func onTextAreaView(service *mockTalkService) view.Middleware {
	return func(uiState state.UIState, context behavior.Context[screen.ViewFunc]) viewmodel.ViewModel {
		text_screen.KeyState.Update(
			uiState.Store,
			context.Target.Name,
			func(s *text_screen.State) {
				s.Buffer = service.text
			},
		)

		return context.Next(uiState)
	}
}

func onKeyEnter(service *mockTalkService) tick.Middleware {
	return func(uiState *state.UIState, event screen.Event, context behavior.Context[screen.TickFunc]) screen.Result {
		currentState, ok := text_screen.KeyState.Get(
			uiState.Store,
			context.Target.Name,
		)

		if !ok || !currentState.Write {
			return context.Next(uiState, event)
		}

		event = screen.NewEvent(
			*key.NewKeyIgnore(),
		)

		result := context.Next(uiState, event)

		text_screen.KeyPulse.Set(
			uiState.Store,
			context.Target.Name,
			true,
		)

		if state, ok := text_screen.KeyState.Get(
			uiState.Store,
			context.Target.Name,
		); ok {
			service.chat = append(service.chat,
				chat.Message{
					Time:    time.Now().UnixMilli(),
					Owner:   service.owner,
					Message: string(state.Buffer),
				},
			)
		}

		service.text = make([]rune, 0)

		return result
	}
}

func wrapStep(vm viewmodel.ViewModel) viewmodel.ViewModel {
	kernel := vm.Kernel.ToUnit()

	pageStep := pageTransformer()

	paddingX := padding.Cols(
		hint.Maximize[winsize.Cols](),
		cols.WithPosition(style.Left),
	)

	paddingY := padding.Rows(
		hint.Maximize[winsize.Rows](),
		rows.WithPosition(style.Top),
	)

	paddingPip := drawable_pipeline.New(kernel).
		SetDrawStep(pageStep).
		PushDataSteps(
			paddingY,
			paddingX,
		).
		ToUnit()

	box := box.New(paddingPip).
		PaddingX(
			hint.Fixed[winsize.Cols](1),
		).
		ToUnit()

	vm.Kernel = stack.NewVStack().
		Push(box)

	return vm
}

func pageTransformer() drawable_pipeline.DrawTransformer {
	action := action.Scroll()
	return func(winsize winsize.Winsize, unit drawable.Unit) ([]text.Line, bool) {
		transformer := focus.DrawTransformer(action)
		return transformer(winsize, unit)
	}
}
