package clip

import (
	"strings"
	"testing"
	"time"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/test/support/mock"

	screen_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/screen"
)

func frameToString(frame *Frame) string {
	var sb strings.Builder
	for i, fs := range frame.frags {
		for _, f := range fs {
			sb.WriteString(f.Text)
		}
		if i < len(frame.frags)-1 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

func TestClip_ToNode(t *testing.T) {
	node := New().
		Name("base").
		ToNode()

	screen_test.Helper_ToNode(t, node)

	assert.Equal(t, node.Name, "base")
}

func TestClip_Stack(t *testing.T) {
	stack := New().ToNode().Stack

	assert.True(t, stack.Has(Name))
}

func TestClip_Boot(t *testing.T) {
	clock := &mock.TestClock{Time: 1000}

	clip := New()

	clip.clock = clock.Now

	clip.frames = []Frame{
		NewFrame(
			TextFrags(),
		),
		NewFrame(
			TextFrags(),
			TextFrags(),
		),
	}

	node := clip.ToNode()

	uiState := state.NewUIState()
	KeyActive.Set(uiState.Store, node.Name, false)

	node.Screen.Boot(*uiState)

	assert.False(t, clip.active)

	assert.Equal(t, clip.start, time.Duration(1000))

	assert.Size(t, 2, clip.empty.frags)
	assert.Size(t, 2, clip.frames[0].frags)
	assert.Size(t, 2, clip.frames[1].frags)
}

func TestClip_MakeEmptyFrame(t *testing.T) {
	frame := makeEmptyFrame(3)

	assert.Size(t, 3, frame.frags)

	for i := range frame.frags {
		assert.Size(t, 1, frame.frags[i])
	}
}

func TestClip_FixFramesSize(t *testing.T) {
	frames := []Frame{
		NewFrame(
			TextFrags(),
		),
		NewFrame(
			TextFrags(),
			TextFrags(),
			TextFrags(),
		),
	}

	result := fixFramesSize(3, frames)

	assert.Size(t, 3, result[0].frags)
	assert.Size(t, 3, result[1].frags)

	for _, row := range result[0].frags {
		assert.Size(t, 1, row)
	}
}

func TestClip_FindFrame_Inactive(t *testing.T) {
	clip := New()

	clip.empty = makeEmptyFrame(1)

	frame, ok := clip.findFrame()

	assert.False(t, ok)
	assert.Size(t, 1, frame.frags)
}

func TestClip_FindFrame_WithoutFrames(t *testing.T) {
	clip := New()

	clip.active = true
	clip.empty = makeEmptyFrame(1)

	frame, ok := clip.findFrame()

	assert.False(t, ok)
	assert.Size(t, 1, frame.frags)
}

func TestClip_FindFrame(t *testing.T) {
	clock := &mock.TestClock{
		Time: 0,
	}

	frame0 := NewFrame(
		TextFrags("0"),
	)

	frame1 := NewFrame(
		TextFrags("1"),
	)

	frame2 := NewFrame(
		TextFrags("2"),
	)

	clip := New().
		SetPause(1000)

	clip.clock = clock.Now
	clip.start = 0

	clip.frames = []Frame{
		frame0,
		frame1,
		frame2,
	}

	frame, ok := clip.findFrame()

	assert.True(t, ok)
	assert.Equal(t, frameToString(&frame0), frameToString(&frame))

	clock.Advance(1000)

	frame, _ = clip.findFrame()
	assert.Equal(t, frameToString(&frame1), frameToString(&frame))

	clock.Advance(1000)

	frame, _ = clip.findFrame()
	assert.Equal(t, frameToString(&frame2), frameToString(&frame))

	clock.Advance(1000)

	frame, _ = clip.findFrame()
	assert.Equal(t, frameToString(&frame0), frameToString(&frame))
}

func TestClip_View_NeedsPulse(t *testing.T) {
	clip := New()

	clip.active = true
	clip.start = 0

	clip.frames = []Frame{
		makeEmptyFrame(1),
	}

	vm := clip.view(*state.NewUIState())

	assert.True(t, vm.Behavior.NeedsPulse)
}

func TestClip_View_DoesNotNeedPulse(t *testing.T) {
	clip := New()

	vm := clip.view(*state.NewUIState())

	assert.False(t, vm.Behavior.NeedsPulse)
}
