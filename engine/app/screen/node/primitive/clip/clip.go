package clip

import (
	"fmt"
	"time"

	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/keymap"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/drain"
	"github.com/Rafael24595/go-reacterm-core/engine/platform/clock"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

const Name = "clip"

const (
	defaultPause = time.Millisecond * 500
	defaultDelta = time.Duration(10)
	defaultLimit = time.Duration(5000)
)

type Clip struct {
	reference  string
	loaded     bool
	clock      clock.Clock
	bindings   *keymap.Bindings[Command]
	definition screen.Definition
	pause      time.Duration
	start      time.Duration
	active     bool
	empty      Frame
	frames     []Frame
}

func New() *Clip {
	return &Clip{
		reference:  Name,
		loaded:     false,
		clock:      clock.UnixMilliClock,
		bindings:   defaultReadBindings,
		definition: screen.EmptyDefinition(),
		pause:      defaultPause,
		start:      0,
		active:     true,
		empty:      NewFrame(),
		frames:     make([]Frame, 0),
	}
}

func (n *Clip) Name(name string) *Clip {
	if n.loaded {
		assert.Unreachable(screen.MessageModified)
		return n
	}

	n.reference = name
	return n
}

func (n *Clip) EnableWriteMode() *Clip {
	if n.loaded {
		assert.Unreachable(screen.MessageModified)
		return n
	}

	n.bindings = defaultWriteBindings
	return n
}

func (n *Clip) WithBindings(overrides *keymap.Bindings[Command]) *Clip {
	if n.loaded {
		assert.Unreachable(screen.MessageModified)
		return n
	}

	n.bindings = n.bindings.Overlay(overrides)
	return n
}

func (n *Clip) Active(active bool) *Clip {
	n.active = active
	return n
}

func (n *Clip) SetPause(pause time.Duration) *Clip {
	if n.loaded {
		assert.Unreachable(screen.MessageModified)
		return n
	}

	n.pause = pause
	return n
}

func (n *Clip) SetFrames(frames ...Frame) *Clip {
	if n.loaded {
		assert.Unreachable(screen.MessageNewElement)
		return n
	}

	n.frames = frames
	return n
}

func (n *Clip) ToNode() screen.Node {
	return screen.NewBuilder().
		Name(n.reference).
		NameToStack().
		Boot(n.boot).
		Keys(n.keys).
		Tick(n.tick).
		View(n.view).
		ToNode()
}

func (n *Clip) boot(uiState state.UIState) {
	if n.loaded {
		return
	}

	n.loadFromStore(uiState)

	n.loaded = true
	n.start = time.Duration(n.clock())

	maxRows := n.maxFrameRows()

	n.empty = makeEmptyFrame(maxRows)
	n.frames = fixFramesSize(maxRows, n.frames)

	n.definition = keymap.BindingsToDefinition(n.bindings)
}

func fixFramesSize(maxRows int, frames []Frame) []Frame {
	fixedFrames := make([]Frame, len(frames))

	for i := range frames {
		fixedFrame := frames[i]
		if len(frames[i].frags) < maxRows {
			fixedFrame.frags = make([]row, maxRows)
			copy(fixedFrame.frags, frames[i].frags)
		}

		fixedFrames[i] = fixedFrame
	}

	return normalizeFrames(fixedFrames...)
}

func makeEmptyFrame(maxRows int) Frame {
	emptyFrame := Frame{
		frags: make([]row, maxRows),
	}

	normalizeFrame(&emptyFrame)
	return emptyFrame
}

func (n *Clip) maxFrameRows() int {
	rows := 0
	for _, f := range n.frames {
		rows = max(rows, len(f.frags))
	}
	return rows
}

func (n *Clip) loadFromStore(uiState state.UIState) {
	sync, _ := KeySync.Take(
		uiState.Store,
		n.reference,
	)

	if sync.Active != nil {
		n.active = *sync.Active
	}

	if sync.Pause != nil {
		n.pause = math.Clamp(
			*sync.Pause, defaultDelta, defaultLimit,
		)
	}

	if restart, ok := KeyRestart.Take(
		uiState.Store,
		n.reference,
	); ok && restart {
		n.start = time.Duration(n.clock())
	}
}

func (n *Clip) keys() screen.Definition {
	return n.definition
}

func (n *Clip) tick(uiState *state.UIState, event screen.Event) screen.Result {
	switch n.bindings.Command(event.Key.Code) {
	case CmdWriteDec:
		n.pause = max(defaultDelta, n.pause-defaultDelta)
	case CmdWriteInc:
		n.pause = min(defaultLimit, n.pause+defaultDelta)
	}

	n.tickToStore(uiState)
	return screen.ResultFromUIState(uiState)
}

func (n *Clip) tickToStore(uiState *state.UIState) {
	state := State{
		Active: n.active,
		Pause:  n.pause,
	}

	KeyState.Set(
		uiState.Store,
		n.reference,
		state,
	)
}

func (n *Clip) view(uiState state.UIState) viewmodel.ViewModel {
	n.loadFromStore(uiState)

	vm := viewmodel.New()

	frame, ok := n.findFrame()
	vm.Behavior.NeedsPulse = ok

	lines := frameToLines(frame)
	vm.Kernel.Push(
		drain.UnitFromLines(lines...),
	)

	if n.bindings.Size() > 0 {
		line := line.New(
			fmt.Sprintf("Speed: %d", n.pause),
		)

		vm.Footer.Push(
			drain.UnitFromLines(*line),
		)
	}

	return *vm
}

func (n *Clip) findFrame() (Frame, bool) {
	framesLen := time.Duration(len(n.frames))
	if !n.active || framesLen == 0 {
		return n.empty, false
	}

	now := time.Duration(n.clock())
	elapsed := now - n.start
	index := (elapsed / n.pause) % framesLen

	return n.frames[index], true
}
