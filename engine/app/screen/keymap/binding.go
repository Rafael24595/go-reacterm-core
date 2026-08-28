package keymap

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/dict"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
)

// Command represents an abstract domain or UI action constraint bound to key inputs.
type Command interface {
	~uint8
}

// Binding associates a specific domain Command with an optional key visual Descriptor.
type Binding[T Command] struct {
	// Command represents the domain command associated with this key binding.
	Command T
	// Descriptor points to the visual key descriptor used for UI help bars or key hints.
	Descriptor *key.Descriptor
}

// Bindings manages a collection of key action bindings mapped to application commands.
type Bindings[T Command] struct {
	keys     *dict.LinkedMap[key.Action, Binding[T]]
	resolver key.DescriptorResolver
}

// NewBindings initializes and returns a new Bindings instance with default configurations.
func NewBindings[T Command]() *Bindings[T] {
	return new(Bindings[T]).lazyInit()
}

func (b *Bindings[T]) lazyInit() *Bindings[T] {
	if b.keys == nil {
		b.keys = dict.NewLinkedMap[key.Action, Binding[T]]()
	}

	if b.resolver == nil {
		b.resolver = key.ResolveDescriptor
	}

	return b
}

// SetResolver sets a custom DescriptorResolver function used to resolve visual descriptors when binding keys.
func (b *Bindings[T]) SetResolver(resolver key.DescriptorResolver) *Bindings[T] {
	b.lazyInit()

	b.resolver = resolver
	return b
}

// Size returns the total number of registered key action bindings.
func (b *Bindings[T]) Size() uint {
	b.lazyInit()

	return b.keys.Size()
}

// Has checks whether a binding exists for the specified key Action.
func (b *Bindings[T]) Has(action key.Action) bool {
	b.lazyInit()

	return b.keys.Exists(action)
}

// Resolve looks up the command associated with the specified key Action, returning a boolean indicating if it was found.
func (b *Bindings[T]) Resolve(action key.Action) (T, bool) {
	b.lazyInit()

	command, ok := b.keys.Get(action)
	if !ok {
		var zero T
		return zero, false
	}

	return command.Command, true
}

// Command retrieves the command associated with the key Action, returning the zero value if not found.
func (b *Bindings[T]) Command(action key.Action) T {
	b.lazyInit()

	command, _ := b.Resolve(action)
	return command
}

// Commands returns a set containing all unique commands registered in the bindings.
func (b *Bindings[T]) Commands() set.Set[T] {
	b.lazyInit()

	commands := set.New[T](int(b.keys.Size()))
	for v := range b.keys.Values() {
		commands.Add(v.Command)
	}
	return commands
}

// Overlay creates a new Bindings instance by merging overrides on top of the current bindings.
func (b *Bindings[T]) Overlay(overrides *Bindings[T]) *Bindings[T] {
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

// Bind registers a command for a key Action and returns the Bindings instance for method chaining.
func (b *Bindings[T]) Bind(
	action key.Action,
	command T,
	descriptors ...key.Descriptor,
) *Bindings[T] {
	b.TryBind(action, command, descriptors...)
	return b
}

// TryBind registers a command for a key Action, returning the previous binding and a boolean indicating replacement.
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

// Clone creates a deep copy of the current Bindings instance.
func (b *Bindings[T]) Clone() *Bindings[T] {
	b.lazyInit()

	result := NewBindings[T]()
	result.keys = b.keys.Clone()
	result.resolver = b.resolver

	return result
}

// BindingsToDefinition converts a Bindings instance into a screen.Definition mapping required keys and descriptors.
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
