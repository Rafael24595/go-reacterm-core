package keymap

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/dict"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
)

type Command interface {
	~uint8
}

type Binding[T Command] struct {
	Command    T
	Descriptor *key.Descriptor
}

type Bindings[T Command] struct {
	keys     *dict.LinkedMap[key.Action, Binding[T]]
	resolver key.DescriptorResolver
}

func NewBindings[T Command]() *Bindings[T] {
	return new(Bindings[T]).lazyInit()
}

func (b *Bindings[T]) lazyInit() *Bindings[T] {
	if b.keys == nil {
		b.keys = dict.NewLinkedMap[key.Action, Binding[T]]()
	}

	if b.resolver == nil {
		b.resolver = key.FindDescriptor
	}

	return b
}

func (b *Bindings[T]) Size() uint {
	b.lazyInit()

	return b.keys.Size()
}

func (b *Bindings[T]) Has(action key.Action) bool {
	b.lazyInit()

	return b.keys.Exists(action)
}

func (b *Bindings[T]) Resolve(action key.Action) (T, bool) {
	b.lazyInit()

	command, ok := b.keys.Get(action)
	if !ok {
		var zero T
		return zero, false
	}

	return command.Command, true
}

func (b *Bindings[T]) Command(action key.Action) T {
	b.lazyInit()

	command, _ := b.Resolve(action)
	return command
}

func (b *Bindings[T]) Commands() set.Set[T] {
	b.lazyInit()

	commands := set.New[T](int(b.keys.Size()))
	for v := range b.keys.Values() {
		commands.Add(v.Command)
	}
	return commands
}

func (b *Bindings[T]) Overlay(
	overrides *Bindings[T],
) *Bindings[T] {
	b.lazyInit()

	result := b.Clone()
	if overrides == nil {
		return result
	}

	result.keys.Merge(
		overrides.keys.Clone(),
	)

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

	previous, replaced := b.keys.Get(action)

	b.keys.Set(action, Binding[T]{
		Command:    command,
		Descriptor: descriptor,
	})

	return previous, replaced
}

func (b *Bindings[T]) Clone() *Bindings[T] {
	b.lazyInit()
	
	result := NewBindings[T]()
	result.keys = b.keys.Clone()
	result.resolver = b.resolver
	return result
}

func BindingsToDefinition[T Command](b *Bindings[T]) screen.Definition {
	required := dict.NewLinkedMap[key.Action, key.Key]()
	descriptor := dict.NewLinkedMap[key.Action, key.Descriptor]()

	for k, v := range b.keys.All() {
		required.Set(k, *key.NewKeyCode(k))
		descriptor.Set(k, *v.Descriptor)
	}

	return screen.Definition{
		RequireKeys: required,
		Descriptor:  descriptor,
	}
}
