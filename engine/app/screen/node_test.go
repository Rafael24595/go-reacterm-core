package screen

import (
	"errors"
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"
)

func TestNode_Compile(t *testing.T) {
	errPassFailed := errors.New("pass error")

	passAddPrefix := func(n Node) (Node, error) {
		n.Name = "Compiled_" + n.Name
		return n, nil
	}

	passWithError := func(n Node) (Node, error) {
		return n, errPassFailed
	}

	tests := []struct {
		name        string
		initialNode Node
		passes      []Pass
		wantName    string
		wantErr     bool
	}{
		{
			name:        "No passes provided",
			initialNode: Node{Name: "Base"},
			passes:      nil,
			wantName:    "Base",
			wantErr:     false,
		},
		{
			name:        "Successful pass execution pipeline",
			initialNode: Node{Name: "Base"},
			passes:      []Pass{passAddPrefix},
			wantName:    "Compiled_Base",
			wantErr:     false,
		},
		{
			name:        "Pipeline failure due to pass error",
			initialNode: Node{Name: "Base"},
			passes:      []Pass{passAddPrefix, passWithError},
			wantName:    "Compiled_Base",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CompileNode(tt.initialNode, tt.passes...)

			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.wantName, got.Name)
		})
	}
}

func TestIsZeroNode(t *testing.T) {
	validScreen := Screen{
		Boot: dummyBoot,
		Keys: dummyKeys,
		Tick: dummyTick,
		View: dummyView,
	}

	validSet := set.New[string]()

	tests := []struct {
		name string
		node Node
		want bool
	}{
		{
			name: "Valid node with all fields initialized",
			node: Node{
				Name:   "ValidNode",
				Tags:   validSet,
				Stack:  validSet,
				Screen: validScreen,
			},
			want: false,
		},
		{
			name: "Empty Name field",
			node: Node{
				Name:   "",
				Tags:   validSet,
				Stack:  validSet,
				Screen: validScreen,
			},
			want: true,
		},
		{
			name: "Nil Tags field",
			node: Node{
				Name:   "Node1",
				Tags:   nil,
				Stack:  validSet,
				Screen: validScreen,
			},
			want: true,
		},
		{
			name: "Nil Stack field",
			node: Node{
				Name:   "Node1",
				Tags:   validSet,
				Stack:  nil,
				Screen: validScreen,
			},
			want: true,
		},
		{
			name: "Screen is a zero screen",
			node: Node{
				Name:   "Node1",
				Tags:   validSet,
				Stack:  validSet,
				Screen: Screen{},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IsZeroNode(tt.node))
		})
	}
}
