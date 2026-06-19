package store

import (
	"fmt"
	"sync"
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"
)

func TestStore_CRUD(t *testing.T) {
	store := New()

	scope := "Landing"
	key := "lang"
	value := "golang"

	store.Push(scope, key, value)
	arg, found := store.Find(scope, key)

	assert.True(t, found)
	assert.Equal(t, value, arg.Stringf())

	store.RemoveArgument(scope, key)
	_, found = store.Find(scope, key)

	assert.False(t, found)

	store.Push(scope, "order", 1)
	store.RemoveScope(scope)
	_, found = store.Find(scope, "order")

	assert.False(t, found)
}

func TestStore_RetainOnly(t *testing.T) {
	store := New()

	store.Push("Home", "a", 1)
	store.Push("Settings", "b", 2)
	store.Push("Profile", "c", 3)

	keep := set.From("Home", "Profile")
	store.RetainOnly(keep)

	_, found := store.Find("Home", "a")
	assert.True(t, found)

	_, found = store.Find("Settings", "b")
	assert.False(t, found)
}

func TestStore_Push_CreatesScope(t *testing.T) {
	store := New()

	store.Push("NewScope", "key", 123)

	_, ok := store.Find("NewScope", "key")
	assert.True(t, ok)
}

func TestStore_Push_Overwrite(t *testing.T) {
	store := New()

	store.Push("S", "k", 1)
	store.Push("S", "k", 2)

	v, ok := store.Find("S", "k")
	assert.True(t, ok)
	assert.Equal(t, 2, v.Intd(0))
}

func TestUpdate_ExistingValue(t *testing.T) {
	store := New()

	var key Key[int] = "counter"

	Push(store, "scope", key, 1)

	value, ok := Update(
		store,
		"scope",
		key,
		func(v *int) {
			*v++
		},
	)

	assert.True(t, ok)
	assert.Equal(t, 2, value)

	found, ok := Find(store, "scope", key)

	assert.True(t, ok)
	assert.Equal(t, 2, found)
}

func TestUpdate_MissingValue(t *testing.T) {
	store := New()

	var key Key[int] = "counter"

	called := 0

	value, ok := Update(
		store,
		"scope",
		key,
		func(v *int) {
			called += 1
			*v += 1
		},
	)

	assert.False(t, ok)
	assert.Equal(t, 0, called)
	assert.Equal(t, 0, value)
}

func TestUpdate_PersistsChanges(t *testing.T) {
	store := New()

	var key Key[string] = "text"

	Push(store, "scope", key, "hello")

	Update(
		store,
		"scope",
		key,
		func(v *string) {
			*v = "world"
		},
	)

	found, ok := Find(store, "scope", key)

	assert.True(t, ok)
	assert.Equal(t, "world", found)
}

func TestUpdate_ReceivesCurrentValue(t *testing.T) {
	store := New()

	var key Key[string] = "text"

	Push(store, "scope", key, "hello")

	Update(
		store,
		"scope",
		key,
		func(v *string) {
			assert.Equal(t, "hello", *v)
		},
	)
}

func TestUpdate_Struct(t *testing.T) {
	type State struct {
		Buffer string
	}

	store := New()

	var key Key[State] = "state"

	Push(
		store,
		"scope",
		key,
		State{Buffer: "old"},
	)

	Update(
		store,
		"scope",
		key,
		func(s *State) {
			s.Buffer = "new"
		},
	)

	found, ok := Find(store, "scope", key)

	assert.True(t, ok)
	assert.Equal(t, "new", found.Buffer)
}

func TestUpsert_InsertWhenMissing(t *testing.T) {
	store := New()

	var key Key[int] = "counter"

	_, ok := Find(
		store,
		"scope",
		key,
	)

	assert.False(t, ok)

	value, ok := Upsert(
		store,
		"scope",
		key,
		func(v *int) {
			*v = 42
		},
	)

	assert.True(t, ok)
	assert.Equal(t, 42, value)

	found, ok := Find(store, "scope", key)

	assert.True(t, ok)
	assert.Equal(t, 42, found)
}

func TestUpsert_UpdateExisting(t *testing.T) {
	store := New()

	var key Key[int] = "counter"

	Push(store, "scope", key, 1)

	value, ok := Find(
		store,
		"scope",
		key,
	)

	assert.True(t, ok)
	assert.Equal(t, 1, value)

	value, ok = Upsert(
		store,
		"scope",
		key,
		func(v *int) {
			*v++
		},
	)

	assert.True(t, ok)
	assert.Equal(t, 2, value)

	found, ok := Find(store, "scope", key)

	assert.True(t, ok)
	assert.Equal(t, 2, found)
}

func TestUpsert_ReceivesCurrentValue(t *testing.T) {
	store := New()

	var key Key[string] = "counter"

	Push(store, "scope", key, "hello")

	Upsert(
		store,
		"scope",
		key,
		func(v *string) {
			*v += " world"
		},
	)

	found, ok := Find(store, "scope", key)

	assert.True(t, ok)
	assert.Equal(t, "hello world", found)
}

func TestUpsert_UsesZeroValueWhenMissing(t *testing.T) {
	store := New()

	var key Key[bool] = "active"

	Upsert(
		store,
		"scope",
		key,
		func(v *bool) {
			assert.False(t, *v)
			*v = true
		},
	)

	found, ok := Find(store, "scope", key)

	assert.True(t, ok)
	assert.True(t, found)
}

func TestUpsert_CreatesEntryEvenWithoutChanges(t *testing.T) {
	store := New()

	var key Key[int] = "counter"

	Upsert(
		store,
		"scope",
		key,
		func(v *int) {},
	)

	found, ok := Find(store, "scope", key)

	assert.True(t, ok)
	assert.Equal(t, 0, found)
}

func TestStore_RemoveScope_DeletesAllKeys(t *testing.T) {
	store := New()

	store.Push("A", "x", 1)
	store.Push("A", "y", 2)

	store.RemoveScope("A")

	_, ok1 := store.Find("A", "x")
	_, ok2 := store.Find("A", "y")

	assert.False(t, ok1)
	assert.False(t, ok2)
}

func TestStore_RetainOnly_EmptySet(t *testing.T) {
	store := New()

	store.Push("A", "x", 1)
	store.Push("B", "y", 2)

	store.RetainOnly(set.Set[string]{})

	_, okA := store.Find("A", "x")
	_, okB := store.Find("B", "y")

	assert.False(t, okA)
	assert.False(t, okB)
}

func TestStore_RetainOnly_NoMatch(t *testing.T) {
	store := New()

	store.Push("A", "x", 1)
	store.Push("B", "y", 2)

	store.RetainOnly(set.From("C"))

	_, okA := store.Find("A", "x")
	_, okB := store.Find("B", "y")

	assert.False(t, okA)
	assert.False(t, okB)
}

func TestStore_LastWriteWins(t *testing.T) {
	store := New()

	store.Push("S", "k", "a")
	store.Push("S", "k", "b")
	store.Push("S", "k", "c")

	v, _ := store.Find("S", "k")

	assert.Equal(t, "c", v.Stringf())
}

func TestStore_Concurrency(t *testing.T) {
	store := New()

	const workers = 15

	var wg sync.WaitGroup
	wg.Add(workers * 2)

	for i := range workers {
		wg.Go(func() {
			defer wg.Done()
			store.Push("Scope", fmt.Sprintf("k%d", i), i)
		})
	}

	for i := range workers {
		wg.Go(func() {
			defer wg.Done()
			store.Find("Scope", fmt.Sprintf("k%d", i))
		})
	}

	wg.Wait()
}

func TestStore_Concurrent_PushSameKey(t *testing.T) {
	store := New()

	const workers = 50

	var wg sync.WaitGroup
	wg.Add(workers)

	for i := range workers {
		wg.Go(func() {
			defer wg.Done()
			store.Push("S", "k", i)
		})
	}

	wg.Wait()

	_, ok := store.Find("S", "k")
	assert.True(t, ok)
}
