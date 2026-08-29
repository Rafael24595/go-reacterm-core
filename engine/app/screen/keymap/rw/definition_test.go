package rw

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
)

func TestEmptyDefinition(t *testing.T) {
	def := EmptyDefinition()

	assert.NotNil(t, def.Read)
	assert.NotNil(t, def.Write)
}

func TestDefinition_Get(t *testing.T) {
	readDef := screen.EmptyDefinition()
	writeDef := screen.EmptyDefinition()

	def := Definition{
		Read:  readDef,
		Write: writeDef,
	}

	gotRead := def.Get(false)
	assert.Equal(t, readDef, gotRead)

	gotWrite := def.Get(true)
	assert.Equal(t, writeDef, gotWrite)
}
