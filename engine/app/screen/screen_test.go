package screen

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
)

func dummyBoot(state.UIState)                     {}
func dummyKeys() Definition                       { return Definition{} }
func dummyTick(*state.UIState, Event) Result      { return Result{} }
func dummyView(state.UIState) viewmodel.ViewModel { return viewmodel.ViewModel{} }

func TestIsZeroScreen(t *testing.T) {

	tests := []struct {
		name     string
		screen   Screen
		expected bool
	}{
		{
			name:     "empty screen (zero value)",
			screen:   Screen{},
			expected: true,
		},
		{
			name: "without Boot",
			screen: Screen{
				Boot: nil,
				Keys: dummyKeys,
				Tick: dummyTick,
				View: dummyView,
			},
			expected: true,
		},
		{
			name: "without Keys",
			screen: Screen{
				Boot: dummyBoot,
				Keys: nil,
				Tick: dummyTick,
				View: dummyView,
			},
			expected: true,
		},
		{
			name: "without Tick",
			screen: Screen{
				Boot: dummyBoot,
				Keys: dummyKeys,
				Tick: nil,
				View: dummyView,
			},
			expected: true,
		},
		{
			name: "without View",
			screen: Screen{
				Boot: dummyBoot,
				Keys: dummyKeys,
				Tick: dummyTick,
				View: nil,
			},
			expected: true,
		},
		{
			name: "all fields defined (valid screen)",
			screen: Screen{
				Boot: dummyBoot,
				Keys: dummyKeys,
				Tick: dummyTick,
				View: dummyView,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsZeroScreen(tt.screen)
			assert.Equal(t, tt.expected, got)
		})
	}
}
