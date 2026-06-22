package keymap

import (
	"maps"

	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
)

type descriptorResolver func(action key.Action) *key.Descriptor

type Command interface {
	~uint8
}

type Binding[T Command] struct {
	Command    T
	Descriptor *key.Descriptor
}

type Bindings[T Command] struct {
	keys     map[key.Action]Binding[T]
	resolver descriptorResolver
}

func NewBindings[T Command]() *Bindings[T] {
	return new(Bindings[T]).lazyInit()
}

func (b *Bindings[T]) lazyInit() *Bindings[T] {
	if b.keys == nil {
		b.keys = make(map[key.Action]Binding[T])
	}

	if b.resolver == nil {
		b.resolver = key.FindDescriptor
	}

	return b
}

func (b *Bindings[T]) Has(action key.Action) bool {
	b.lazyInit()

	_, ok := b.keys[action]
	return ok
}

func (b *Bindings[T]) Resolve(action key.Action) (T, bool) {
	b.lazyInit()

	command, ok := b.keys[action]
	if !ok {
		var zero T
		return zero, false
	}

	return command.Command, true
}

func (b *Bindings[T]) Overlay(
	overrides *Bindings[T],
) *Bindings[T] {
	b.lazyInit()

	result := b.Clone()
	if overrides == nil {
		return result
	}

	maps.Copy(result.keys, overrides.keys)
	return result
}

func (b *Bindings[T]) Bind(
	action key.Action,
	command T,
	descriptors ...key.Descriptor,
) *Bindings[T] {
	b.TryBind(action, command, descriptors...)
	return b
}

func (b *Bindings[T]) TryBind(
	action key.Action,
	command T,
	descriptors ...key.Descriptor,
) (Binding[T], bool) {
	b.lazyInit()

	var descriptor *key.Descriptor
	if len(descriptors) > 0 {
		descriptor = &descriptors[0]
	} else {
		descriptor = b.resolver(action)
	}

	previous, replaced := b.keys[action]

	b.keys[action] = Binding[T]{
		Command:    command,
		Descriptor: descriptor,
	}

	return previous, replaced
}

func (b *Bindings[T]) Clone() *Bindings[T] {
	b.lazyInit()
	
	result := NewBindings[T]()
	maps.Copy(result.keys, b.keys)
	result.resolver = b.resolver
	return result
}
