package set

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestSet_Has(t *testing.T) {
	s := From("apple", "banana")

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Element exists",
			input:    "apple",
			expected: true,
		},
		{
			input:    "Element does not exist",
			name:     "orange",
			expected: false,
		},
		{
			input:    "Empty string",
			name:     "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, s.Has(tt.input))
		})
	}
}

func TestSet_Any(t *testing.T) {
	tests := []struct {
		name     string
		setA     []int
		setB     []int
		expected bool
	}{
		{
			name:     "Direct intersection",
			setA:     []int{1, 2, 3},
			setB:     []int{3, 4, 5},
			expected: true,
		},
		{
			name:     "No intersection",
			setA:     []int{1, 2},
			setB:     []int{3, 4},
			expected: false,
		},
		{
			name:     "One set is empty",
			setA:     []int{1, 2},
			setB:     []int{},
			expected: false,
		},
		{
			name:     "Both sets are empty",
			setA:     []int{},
			setB:     []int{},
			expected: false,
		},
		{
			name:     "Identical sets",
			setA:     []int{10, 20},
			setB:     []int{10, 20},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := From(tt.setA...)
			b := From(tt.setB...)

			assert.Equal(t, tt.expected, a.Any(b))
			assert.Equal(t, tt.expected, b.Any(a))
		})
	}
}

func TestSet_Add(t *testing.T) {
	s := New[int](1)
	s.Add(42)

	assert.Len(t, 1, s)
	assert.True(t, s.Has(42))
}

func TestSet_Merge(t *testing.T) {
	s1 := From(1, 2, 3)
	s2 := From(3, 4, 5)

	s1.Merge(s2)

	assert.Len(t, 5, s1)

	assert.True(t, s1.Has(4))
	assert.True(t, s1.Has(5))
}

func TestSet_Clone(t *testing.T) {
	s1 := From(1, 2, 3)
	s2 := s1.Clone()

	s2.Add(4, 5, 6)

	assert.Len(t, 3, s1)
	assert.Len(t, 6, s2)

	assert.False(t, s1.Has(6))
	assert.True(t, s2.Has(6))
}

func BenchmarkSet_Any(b *testing.B) {
	large := New[int](1000)
	for i := range 1000 {
		large.Add(i)
	}

	small := New[int](2)
	small.Add(999)
	small.Add(2000)

	for b.Loop() {
		large.Any(small)
	}
}
