package screen

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/dict"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
)

func TestEmptyDefinition(t *testing.T) {
	def := EmptyDefinition()

	assert.NotNil(t, def.RequireKeys)
	assert.NotNil(t, def.Descriptor)
	assert.Equal(t, 0, def.RequireKeys.Size())
	assert.Equal(t, 0, def.Descriptor.Size())
}

func TestNewDefinition(t *testing.T) {
	mockResolver := func(actions ...key.Action) *dict.LinkedMap[key.Action, key.Descriptor] {
		m := dict.NewLinkedMap[key.Action, key.Descriptor]()
		for _, act := range actions {
			m.Set(act, key.Descriptor{})
		}
		return m
	}

	actSubmit := key.Action(1)
	actCancel := key.Action(2)

	def := NewDefinition(mockResolver, actSubmit, actCancel)

	assert.Equal(t, 2, def.RequireKeys.Size())
	assert.True(t, def.RequireKeys.Exists(actSubmit))
	assert.True(t, def.RequireKeys.Exists(actCancel))
	assert.Equal(t, 2, def.Descriptor.Size())
}

func TestDefinition_IsRequired(t *testing.T) {
	actSubmit := key.Action(1)
	actCancel := key.Action(2)

	t.Run("Action is explicitly required", func(t *testing.T) {
		def := EmptyDefinition()
		def.RequireKeys.Set(actSubmit, key.Key{Code: actSubmit})

		assert.True(t, def.IsRequired(key.Key{Code: actSubmit}))
		assert.False(t, def.IsRequired(key.Key{Code: actCancel}))
	})

	t.Run("ActionAll wildcard marks any key as required", func(t *testing.T) {
		def := EmptyDefinition()
		def.RequireKeys.Set(key.ActionAll, key.Key{Code: key.ActionAll})

		assert.True(t, def.IsRequired(key.Key{Code: actSubmit}))
		assert.True(t, def.IsRequired(key.Key{Code: actCancel}))
	})
}

func TestDefinition_Merge(t *testing.T) {
	act1 := key.Action(1)
	act2 := key.Action(2)

	def1 := EmptyDefinition()
	def1.RequireKeys.Set(act1, key.Key{Code: act1})
	def1.Descriptor.Set(act1, key.Descriptor{})

	def2 := EmptyDefinition()
	def2.RequireKeys.Set(act2, key.Key{Code: act2})
	def2.Descriptor.Set(act2, key.Descriptor{})

	merged := def1.Merge(def2)

	assert.Equal(t, 2, merged.RequireKeys.Size())
	assert.True(t, merged.RequireKeys.Exists(act1))
	assert.True(t, merged.RequireKeys.Exists(act2))

	assert.Equal(t, 2, merged.Descriptor.Size())
}