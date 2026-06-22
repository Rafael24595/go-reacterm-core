package keymap

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
)

type testCommand uint8

const (
	CmdNone testCommand = iota
	CmdOpen
	CmdClose
	CmdSave
)

func mockResolver(key.Action) *key.Descriptor {
	return nil
}

func TestNewKeysBindings(t *testing.T) {
	kb := NewBindings[testCommand]()

	assert.NotNil(t, kb)
	assert.NotNil(t, kb.keys)
}

func TestLazyInit(t *testing.T) {
	var kb Bindings[testCommand]

	assert.False(t, kb.Has(key.Action(1)))

	kb.TryBind(key.Action(1), CmdOpen)

	assert.NotNil(t, kb.keys)
}

func TestBindAndHas(t *testing.T) {
	kb := NewBindings[testCommand]()
	action := key.Action(10)

	assert.False(t, kb.Has(action))

	kb.Bind(action, CmdOpen)

	assert.True(t, kb.Has(action))
}

func TestResolve(t *testing.T) {
	kb := NewBindings[testCommand]()
	action := key.Action(20)

	cmd, ok := kb.Resolve(action)

	assert.False(t, ok)
	assert.Equal(t, CmdNone, cmd)

	kb.Bind(action, CmdSave)
	cmd, ok = kb.Resolve(action)

	assert.True(t, ok)
	assert.Equal(t, CmdSave, cmd)
}

func TestTryBind(t *testing.T) {
	kb := NewBindings[testCommand]()
	kb.resolver = mockResolver

	action := key.Action(30)

	_, replaced := kb.TryBind(action, CmdOpen)

	assert.False(t, replaced)

	prev, replaced := kb.TryBind(action, CmdClose)

	assert.True(t, replaced)
	assert.Equal(t, CmdOpen, prev.Command)

	cmd, _ := kb.Resolve(action)
	assert.Equal(t, CmdClose, cmd)
}

func TestClone(t *testing.T) {
	kb1 := NewBindings[testCommand]()
	kb1.resolver = mockResolver

	action := key.Action(40)

	kb1.Bind(action, CmdOpen)
	kb2 := kb1.Clone()

	cmd, ok := kb2.Resolve(action)

	assert.True(t, ok)
	assert.Equal(t, CmdOpen, cmd)

	kb2.Bind(action, CmdClose)
	cmdOrig, _ := kb1.Resolve(action)

	assert.Equal(t, CmdOpen, cmdOrig)
}

func TestOverlay(t *testing.T) {
	kbBase := NewBindings[testCommand]()
	kbBase.Bind(key.Action(1), CmdOpen)
	kbBase.Bind(key.Action(2), CmdClose)

	kbOverrides := NewBindings[testCommand]()
	kbOverrides.Bind(key.Action(2), CmdSave)
	kbOverrides.Bind(key.Action(3), CmdOpen)

	kbNilResult := kbBase.Overlay(nil)

	assert.NotEqual(t, kbBase, kbNilResult)

	kbResult := kbBase.Overlay(kbOverrides)

	cmdBase, _ := kbBase.Resolve(key.Action(2))

	assert.Equal(t, CmdClose, cmdBase)

	cmd, ok := kbResult.Resolve(key.Action(1))

	assert.True(t, ok)
	assert.Equal(t, CmdOpen, cmd)

	cmd, ok = kbResult.Resolve(key.Action(2))

	assert.True(t, ok)
	assert.Equal(t, CmdSave, cmd)

	cmd, ok = kbResult.Resolve(key.Action(3))

	assert.True(t, ok)
	assert.Equal(t, CmdOpen, cmd)
}
