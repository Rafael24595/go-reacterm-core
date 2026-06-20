package clip

import (
	"time"

	assert "github.com/Rafael24595/go-assert/assert/runtime"
	
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/drain"
	"github.com/Rafael24595/go-reacterm-core/engine/platform/clock"
)

const Name = "clip"

const defaultPause = time.Millisecond * 500

type Clip struct {
	reference string
	loaded    bool
	clock     clock.Clock
	pause     time.Duration
	start     time.Duration
	active    bool
	empty     Frame
	frames    []Frame
}

func New() *Clip {
	return &Clip{
		reference: Name,
		loaded:    false,
		clock:     clock.UnixMilliClock,
		pause:     defaultPause,
		start:     0,
		active:    true,
		empty:     NewFrame(),
		frames:    make([]Frame, 0),
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
		WithoutKeys().
		WithoutTick().
		Boot(n.boot).
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
	if active, ok := KeyActive.Get(
		uiState.Store,
		n.reference,
	); ok {
		n.active = active
	}

	if restart, ok := KeyRestart.Get(
		uiState.Store,
		n.reference,
	); ok && restart {
		n.start = time.Duration(n.clock())
	}
}

func (n *Clip) view(state.UIState) viewmodel.ViewModel {
	vm := viewmodel.New()

	frame, ok := n.findFrame()
	vm.Behavior.NeedsPulse = ok

	line := frameToLines(frame)
	vm.Kernel.Push(
		drain.UnitFromLines(line...),
	)

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
