package slot

import (
	"sync"
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestTakeEmpty(t *testing.T) {
	s := New[int]()

	v, ok := s.Take()

	assert.False(t, ok)
	assert.Equal(t, 0, v)
}

func TestSetAndTake(t *testing.T) {
	s := New[string]()

	s.Set("golang")

	v, ok := s.Take()

	assert.True(t, ok)
	assert.Equal(t, "golang", v)
}

func TestTakeConsumesValue(t *testing.T) {
	s := New[string]()

	s.Set("golang")

	v, ok := s.Take()
	assert.True(t, ok)
	assert.Equal(t, "golang", v)

	v, ok = s.Take()
	assert.False(t, ok)
	assert.Equal(t, "", v)
}

func TestSetOverridesPreviousValue(t *testing.T) {
	s := New[int]()

	s.Set(1)
	s.Set(2)

	v, ok := s.Take()

	assert.True(t, ok)
	assert.Equal(t, 2, v)
}

func TestReuseAfterTake(t *testing.T) {
	s := New[int]()

	s.Set(1)

	v, ok := s.Take()
	assert.True(t, ok)
	assert.Equal(t, 1, v)

	s.Set(2)

	v, ok = s.Take()
	assert.True(t, ok)
	assert.Equal(t, 2, v)
}

func TestConcurrentAccess(t *testing.T) {
	s := New[int]()

	var wg sync.WaitGroup

	for i := range 1000 {
		wg.Go(func() {
			s.Set(i)
		})

		wg.Go(func() {
			s.Take()
		})
	}

	wg.Wait()
}
