package store

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

type Lang struct {
	Name    string
	Creator string
}

func TestKey_TypeAndCode(t *testing.T) {
	var key Key[Lang] = "golang"

	assert.Equal(t, "golang", key.Code())
	assert.Equal(t, Lang{}, key.Type())
}

func TestKey_FluentAPI(t *testing.T) {
	const scope = "Langs"
	var currentLang Key[Lang] = "current_lang"

	t.Run("Set and Get", func(t *testing.T) {
		store := New()

		currentLang.Set(store, scope, Lang{
			Name:    "golang",
			Creator: "Ken Thompson",
		})

		user, found := currentLang.Get(store, scope)
		assert.True(t, found)
		assert.Equal(t, "golang", user.Name)
		assert.Equal(t, "Ken Thompson", user.Creator)
	})

	t.Run("Update existing key", func(t *testing.T) {
		store := New()

		currentLang.Set(store, scope, Lang{
			Name:    "golang",
			Creator: "Ken Thompson",
		})

		currentLang.Update(store, scope, func(u *Lang) {
			u.Creator = "Robert Griesemer;Rob Pike;Ken Thompson"
		})

		updatedUser, found := currentLang.Get(store, scope)
		assert.True(t, found)
		assert.Equal(t, "Robert Griesemer;Rob Pike;Ken Thompson", updatedUser.Creator)
	})

	t.Run("Upsert on non-existing key", func(t *testing.T) {
		store := New()
		var otherLang Key[Lang] = "other_lang"

		otherLang.Upsert(store, scope, func(u *Lang) {
			u.Name = "zig"
		})

		cfg, found := otherLang.Get(store, scope)
		assert.True(t, found)
		assert.Equal(t, "zig", cfg.Name)
	})

	t.Run("Take removes and returns value", func(t *testing.T) {
		store := New()

		currentLang.Set(store, scope, Lang{
			Name:    "zig",
			Creator: "Andrew Kelley",
		})

		takenLang, okTake := currentLang.Take(store, scope)
		assert.True(t, okTake)
		assert.Equal(t, "zig", takenLang.Name)

		_, stillExists := currentLang.Get(store, scope)
		assert.False(t, stillExists)
	})

	t.Run("Delete removes value silently", func(t *testing.T) {
		store := New()

		currentLang.Set(store, scope, Lang{
			Name:    "zig",
			Creator: "Andrew Kelley",
		})

		currentLang.Delete(store, scope)

		_, exists := currentLang.Get(store, scope)
		assert.False(t, exists)
	})
}
