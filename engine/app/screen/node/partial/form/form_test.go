package form

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/config/entry"

	screen_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/screen"
)

func TestForm_ToNode(t *testing.T) {
	node := New().ToNode()
	screen_test.Helper_ToNode(t, node)

	assert.Equal(t, Name, node.Name)
}

func TestForm_Propagate(t *testing.T) {
	name := "base"

	node := New().
		AddNode(
			screen_test.MockByName(name),
			entry.Selectable(),
		).
		ToNode()

	screen_test.Helper_Propagate(t, name, 0, node)
}

func TestSetCursor(t *testing.T) {
	tests := []struct {
		name  string
		items []entry.Entry
		input uint16
		want  uint16
	}{
		{
			name:  "empty form returns 0",
			items: make([]entry.Entry, 0),
			input: 5,
			want:  0,
		},
		{
			name:  "within limits",
			items: make([]entry.Entry, 3),
			input: 1,
			want:  1,
		},
		{
			name:  "at upper limit",
			items: make([]entry.Entry, 3),
			input: 2,
			want:  2,
		},
		{
			name:  "outside limits (len) truncates to last index",
			items: make([]entry.Entry, 3),
			input: 3,
			want:  2,
		},
		{
			name:  "outside limits (large) truncates to last index",
			items: make([]entry.Entry, 5),
			input: 999,
			want:  4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Form{
				items: tt.items,
			}

			got := f.setCursor(tt.input)

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want, f.cursor)
		})
	}
}

func TestIncCursor(t *testing.T) {
	tests := []struct {
		name   string
		cursor uint16
		items  []entry.Entry
		step   uint16
		want   uint16
	}{
		{
			name:   "next selectable",
			cursor: 0,
			items: []entry.Entry{
				{Selectable: true},
				{Selectable: true},
				{Selectable: true},
			},
			step: 1,
			want: 1,
		},
		{
			name:   "skip non selectable",
			cursor: 0,
			items: []entry.Entry{
				{Selectable: true},
				{Selectable: false},
				{Selectable: false},
				{Selectable: true},
			},
			step: 1,
			want: 3,
		},
		{
			name:   "wrap around",
			cursor: 3,
			items: []entry.Entry{
				{Selectable: true},
				{Selectable: true},
				{Selectable: true},
				{Selectable: true},
			},
			step: 1,
			want: 0,
		},
		{
			name:   "wrap and skip non selectable",
			cursor: 3,
			items: []entry.Entry{
				{Selectable: false},
				{Selectable: true},
				{Selectable: false},
				{Selectable: true},
			},
			step: 1,
			want: 1,
		},
		{
			name:   "large step",
			cursor: 0,
			items: []entry.Entry{
				{Selectable: true},
				{Selectable: true},
				{Selectable: true},
				{Selectable: true},
			},
			step: 5,
			want: 1,
		},
		{
			name:   "no selectable",
			cursor: 2,
			items: []entry.Entry{
				{Selectable: false},
				{Selectable: false},
				{Selectable: false},
			},
			step: 1,
			want: 2,
		},
		{
			name:   "search starts at offset",
			cursor: 0,
			items: []entry.Entry{
				{Selectable: true},
				{Selectable: false},
				{Selectable: false},
				{Selectable: true},
				{Selectable: true},
			},
			step: 3,
			want: 3,
		},
		{
			name:   "first selectable",
			cursor: 0,
			items: []entry.Entry{
				{Selectable: false},
				{Selectable: true},
				{Selectable: false},
				{Selectable: true},
				{Selectable: true},
			},
			step: 0,
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Form{
				cursor: tt.cursor,
				items:  tt.items,
			}

			got := f.incCursor(tt.step)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDecCursor(t *testing.T) {
	tests := []struct {
		name   string
		cursor uint16
		items  []entry.Entry
		step   uint16
		want   uint16
	}{
		{
			name:   "previous selectable",
			cursor: 2,
			items: []entry.Entry{
				{Selectable: true},
				{Selectable: true},
				{Selectable: true},
			},
			step: 1,
			want: 1,
		},
		{
			name:   "skip non selectable",
			cursor: 3,
			items: []entry.Entry{
				{Selectable: true},
				{Selectable: false},
				{Selectable: false},
				{Selectable: true},
			},
			step: 1,
			want: 0,
		},
		{
			name:   "wrap around",
			cursor: 0,
			items: []entry.Entry{
				{Selectable: true},
				{Selectable: true},
				{Selectable: true},
				{Selectable: true},
			},
			step: 1,
			want: 3,
		},
		{
			name:   "wrap and skip non selectable",
			cursor: 1,
			items: []entry.Entry{
				{Selectable: true},
				{Selectable: false},
				{Selectable: false},
				{Selectable: true},
			},
			step: 2,
			want: 3,
		},
		{
			name:   "large step",
			cursor: 0,
			items: []entry.Entry{
				{Selectable: true},
				{Selectable: true},
				{Selectable: true},
				{Selectable: true},
			},
			step: 5,
			want: 3,
		},
		{
			name:   "last selectable",
			cursor: 4,
			items: []entry.Entry{
				{Selectable: true},
				{Selectable: true},
				{Selectable: false},
				{Selectable: true},
				{Selectable: false},
			},
			step: 0,
			want: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Form{
				cursor: tt.cursor,
				items:  tt.items,
			}

			got := f.decCursor(tt.step)

			assert.Equal(t, tt.want, got)
		})
	}
}
