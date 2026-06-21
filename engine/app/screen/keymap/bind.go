package keymap

import (
	"maps"

	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
)

type descriptorResolver func(action key.Action) *key.Descriptor

type Command interface {
	~uint8
}

type KeyBinding[T Command] struct {
	Command    T
	Descriptor *key.Descriptor
}

type KeysBindings[T Command] struct {
	keys     map[key.Action]KeyBinding[T]
	resolver descriptorResolver
}

func NewKeysBindings[T Command]() *KeysBindings[T] {
	return new(KeysBindings[T]).lazyInit()
}

func (b *KeysBindings[T]) lazyInit() *KeysBindings[T] {
	if b.keys == nil {
		b.keys = make(map[key.Action]KeyBinding[T])
	}

	if b.resolver == nil {
		b.resolver = key.FindDescriptor
	}

	return b
}

func (b *KeysBindings[T]) Has(action key.Action) bool {
	_, ok := b.keys[action]
	return ok
}

func (b *KeysBindings[T]) Resolve(action key.Action) (T, bool) {
	command, ok := b.keys[action]
	if !ok {
		var zero T
		return zero, false
	}
	return command.Command, true
}

func (b *KeysBindings[T]) Overlay(
	overrides *KeysBindings[T],
) *KeysBindings[T] {
	result := b.Clone()
	if overrides == nil {
		return result
	}

	maps.Copy(result.keys, overrides.keys)
	return result
}

func (b *KeysBindings[T]) Bind(
	action key.Action,
	command T,
	descriptors ...key.Descriptor,
) *KeysBindings[T] {
	b.TryBind(action, command, descriptors...)
	return b
}

func (b *KeysBindings[T]) TryBind(
	action key.Action,
	command T,
	descriptors ...key.Descriptor,
) (KeyBinding[T], bool) {
	b.lazyInit()

	var descriptor *key.Descriptor
	if len(descriptors) > 0 {
		descriptor = &descriptors[0]
	} else {
		descriptor = b.resolver(action)
	}

	previous, replaced := b.keys[action]

	b.keys[action] = KeyBinding[T]{
		Command:    command,
		Descriptor: descriptor,
	}

	return previous, replaced
}

func (b *KeysBindings[T]) Clone() *KeysBindings[T] {
	result := NewKeysBindings[T]()
	maps.Copy(result.keys, b.keys)
	result.resolver = b.resolver
	return result
}
