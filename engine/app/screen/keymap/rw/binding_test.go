package rw

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/keymap"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
)

type readCmd uint8
type writeCmd uint8

const (
	CmdReadView  readCmd  = iota + 1
	CmdWriteSave writeCmd = iota + 1
)

func mockResolver(key.Action) *key.Descriptor {
	return &key.Descriptor{}
}

func TestDefinitionFromBindings(t *testing.T) {
	readBindings := keymap.NewBindings[readCmd]().SetResolver(mockResolver)
	readBindings.Bind(10, CmdReadView)

	writeBindings := keymap.NewBindings[writeCmd]().SetResolver(mockResolver)
	writeBindings.Bind(20, CmdWriteSave)

	rwBindings := Bindings[readCmd, writeCmd]{
		Read:  readBindings,
		Write: writeBindings,
	}

	def := DefinitionFromBindings(rwBindings)

	assert.NotNil(t, def.Read.RequireKeys)
	assert.NotNil(t, def.Write.RequireKeys)

	assert.True(t, def.Read.RequireKeys.Exists(key.Action(10)))
	assert.False(t, def.Read.RequireKeys.Exists(key.Action(20)))

	assert.True(t, def.Write.RequireKeys.Exists(key.Action(20)))
	assert.False(t, def.Write.RequireKeys.Exists(key.Action(10)))
}
