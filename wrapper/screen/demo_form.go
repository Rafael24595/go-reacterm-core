package wrapper_screen

import (
	"math/rand"
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
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/primitive/clip"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/primitive/talk"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/config/entry"
	"github.com/Rafael24595/go-reacterm-core/engine/config/layer"
	"github.com/Rafael24595/go-reacterm-core/engine/config/padding/cols"
	"github.com/Rafael24595/go-reacterm-core/engine/config/padding/rows"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
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

type eventKind int

const (
	eventTypingStarted eventKind = iota
	eventTypingStopped
)

type clock[T math.Number] struct {
	rng *rand.Rand
}

func newClock[T math.Number](source int64) *clock[T] {
	return &clock[T]{
		rng: rand.New(
			rand.NewSource(source),
		),
	}
}

func (c *clock[T]) rand(n T) T {
	return T(c.rng.Intn(int(n)))
}

type mockMessageService struct {
	clock    *clock[int]
	messages []chat.Message
}

func (s *mockMessageService) randMessage() chat.Message {
	messagesLen := len(s.messages)
	if messagesLen == 0 {
		return chat.Message{
			Time:    time.Now().UnixMilli(),
			Owner:   "System",
			Message: "Empty messages buffer",
		}
	}

	index := s.clock.rng.Intn(messagesLen)
	message := s.messages[index]
	message.Time = time.Now().UnixMilli()

	return message
}

type mockTalkService struct {
	clock    *clock[time.Duration]
	messages *mockMessageService
	event    eventKind
	owner    string
	text     []rune
	chat     []chat.Message
}

func (s *mockTalkService) reply() {
	go func() {
		time.Sleep(
			1000 + time.Millisecond*s.clock.rand(2000),
		)

		s.event = eventTypingStarted

		time.Sleep(
			1000 + time.Millisecond*s.clock.rand(9000),
		)

		s.event = eventTypingStopped

		s.chat = append(s.chat, s.messages.randMessage())
	}()
}

func NewDemoForm() screen.Node {
	service := initTalkService()

	node := form.New().
		AddNode(
			makeTalk(service),
			entry.Selectable(),
			entry.WithLayout(
				layer.Percent[winsize.Rows](73),
			),
		).
		AddBreak(1).
		AddNode(
			makeClip(service),
			entry.WithLayout(
				layer.Static[winsize.Rows](),
			),
		).
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

func initTalkService() *mockTalkService {
	return &mockTalkService{
		clock:    newClock[time.Duration](12345),
		messages: initMessagesService(),
		event:    eventTypingStopped,
		owner:    "human_001",
		text:     []rune("Hello Golang!"),
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

func initMessagesService() *mockMessageService {
	r := rand.New(rand.NewSource(12345))

	users := []string{
		"human_002",
		"human_003",
		"human_004",
		"human_005",
		"human_006",
	}

	texts := []string{
		"Lorem ipsum dolor sit amet.",
		"Consectetur adipiscing elit.",
		"Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
		"Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.",
		"Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.",
		"Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
		"Integer posuere erat a ante venenatis dapibus posuere velit aliquet.",
		"Curabitur blandit tempus porttitor.",
		"Praesent commodo cursus magna, vel scelerisque nisl consectetur et.",
		"Donec ullamcorper nulla non metus auctor fringilla.",
		"Vestibulum id ligula porta felis euismod semper.",
		"Maecenas faucibus mollis interdum.",
		"Morbi leo risus, porta ac consectetur ac, vestibulum at eros.",
		"Aenean lacinia bibendum nulla sed consectetur.",
		"Etiam porta sem malesuada magna mollis euismod.",
		"Nullam quis risus eget urna mollis ornare vel eu leo.",
		"Cras mattis consectetur purus sit amet fermentum.",
		"Sed posuere consectetur est at lobortis.",
		"Fusce dapibus, tellus ac cursus commodo, tortor mauris condimentum nibh.",
		"Vivamus sagittis lacus vel augue laoreet rutrum faucibus dolor auctor.",
		"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Integer posuere erat a ante venenatis dapibus posuere velit aliquet.",
		"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec sed odio dui. Nulla vitae elit libero, a pharetra augue.",
		"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum id ligula porta felis euismod semper. Cras mattis consectetur purus sit amet fermentum.",
		"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed posuere consectetur est at lobortis. Aenean lacinia bibendum nulla sed consectetur. Donec ullamcorper nulla non metus auctor fringilla.",
		"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Praesent commodo cursus magna, vel scelerisque nisl consectetur et. Curabitur blandit tempus porttitor. Integer posuere erat a ante venenatis dapibus posuere velit aliquet. Maecenas faucibus mollis interdum.",
	}

	messages := make([]chat.Message, 500)

	for i := range messages {
		messages[i] = chat.Message{
			Time:    0,
			Owner:   users[r.Intn(len(users))],
			Message: texts[r.Intn(len(texts))],
		}
	}

	return &mockMessageService{
		clock:    newClock[int](12345),
		messages: messages,
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
		AddText(string(service.text)).
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

		service.reply()

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

func makeClip(service *mockTalkService) screen.Node {
	clip := clip.New().
		Name("article - dolor").
		Active(false).
		SetPause(150).
		SetFrames(
			clip.NewFrame(
				clip.TextFrags(".", ".", "."),
			),
			clip.NewFrame(
				clip.TextFrags("·", ".", "."),
			),
			clip.NewFrame(
				clip.TextFrags("˙", ".", "."),
			),
			clip.NewFrame(
				clip.TextFrags("·", "·", "."),
			),
			clip.NewFrame(
				clip.TextFrags(".", "˙", "."),
			),
			clip.NewFrame(
				clip.TextFrags(".", "·", "·"),
			),
			clip.NewFrame(
				clip.TextFrags(".", ".", "˙"),
			),
			clip.NewFrame(
				clip.TextFrags(".", ".", "·"),
			),
			clip.NewFrame(
				clip.TextFrags(".", ".", "."),
			),
		).
		ToNode()

	clip = view.Use(
		clip,
		onClipView(service),
	)

	return clip
}

func onClipView(service *mockTalkService) view.Middleware {
	return func(uiState state.UIState, context behavior.Context[screen.ViewFunc]) viewmodel.ViewModel {
		show := service.event == eventTypingStarted
		restart := false

		clip.KeyActive.Upsert(
			uiState.Store,
			context.Target.Name,
			func(b *bool) {
				if !*b {
					restart = true
				}
				*b = show
			},
		)

		clip.KeyRestart.Set(
			uiState.Store,
			context.Target.Name,
			restart,
		)

		return context.Next(uiState)
	}
}
