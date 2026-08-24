package screen

import (
	"strconv"

	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"
	"github.com/Rafael24595/go-reacterm-core/engine/platform/clock"
)

const (
	// ErrorMissingName indicates that the Node name is empty.
	ErrorMissingName = "missing_name"
	// ErrorMissingBoot indicates that the Boot function is not set.
	ErrorMissingBoot = "missing_boot"
	// ErrorMissingKeys indicates that the Keys function is not set.
	ErrorMissingKeys = "missing_keys"
	// ErrorMissingTick indicates that the Tick function is not set.
	ErrorMissingTick = "missing_tick"
	// ErrorMissingView indicates that the View function is not set.
	ErrorMissingView = "missing_view"
)

func withoutBoot(state.UIState) {}

func withoutKeys() Definition {
	return EmptyDefinition()
}

func withoutTick(*state.UIState, Event) Result {
	return EmptyResult()
}

// Builder facilitates fluent construction of Node instances and their underlying Screen component.
type Builder struct {
	clock    clock.Clock
	name     string
	stack    set.Set[string]
	children []Node
	boot     BootFunc
	keys     KeysFunc
	tick     TickFunc
	view     ViewFunc
}

// NewBuilder initializes a new Builder instance with default global configurations.
func NewBuilder() *Builder {
	return &Builder{
		clock:    clock.GlobalCounterClock,
		name:     "",
		stack:    set.New[string](),
		children: make([]Node, 0),
		boot:     nil,
		keys:     nil,
		tick:     nil,
		view:     nil,
	}
}

// WithClock overrides the default clock generator used for ID generation.
func (b *Builder) WithClock(clock clock.Clock) *Builder {
	if clock == nil {
		return b
	}

	b.clock = clock
	return b
}

// WithNode copies both metadata and screen functions from an existing Node into the Builder.
func (b *Builder) WithNode(node Node) *Builder {
	return b.WithNodeMeta(node).
		WithNodeScreen(node)
}

// WithNodeMeta copies the name and stack hierarchy from an existing Node into the Builder.
func (b *Builder) WithNodeMeta(node Node) *Builder {
	return b.Name(node.Name).
		Stack(node.Stack)
}

// WithNodeScreen copies the lifecycle functions (Boot, Keys, Tick, View) from an existing Node.
func (b *Builder) WithNodeScreen(node Node) *Builder {
	return b.Boot(node.Screen.Boot).
		Keys(node.Screen.Keys).
		Tick(node.Screen.Tick).
		View(node.Screen.View)
}

// Name assigns the user-defined name for the Node.
func (b *Builder) Name(name string) *Builder {
	b.name = name
	return b
}

// NameToStack adds the current Node name to its stack hierarchy set.
func (b *Builder) NameToStack() *Builder {
	return b.Stack(
		set.From(b.name),
	)
}

// NameAsStack sets the Node name and immediately registers it into the stack hierarchy.
func (b *Builder) NameAsStack(name string) *Builder {
	return b.Name(name).
		NameToStack()
}

// Stack merges a set of stack elements into the Node's existing stack configuration.
func (b *Builder) Stack(stack set.Set[string]) *Builder {
	b.stack.Merge(stack)
	return b
}

// Children appends child nodes to the Node being built.
func (b *Builder) Children(children ...Node) *Builder {
	b.children = append(b.children, children...)
	return b
}

// Boot sets the BootFunc lifecycle callback.
func (b *Builder) Boot(boot BootFunc) *Builder {
	b.boot = boot
	return b
}

// WithoutBoot assigns a default no-op BootFunc lifecycle callback.
func (b *Builder) WithoutBoot() *Builder {
	b.boot = withoutBoot
	return b
}

// Keys sets the KeysFunc key bindings callback.
func (b *Builder) Keys(keys KeysFunc) *Builder {
	b.keys = keys
	return b
}

// WithoutKeys assigns a default empty KeysFunc callback returning no shortcuts.
func (b *Builder) WithoutKeys() *Builder {
	b.keys = withoutKeys
	return b
}

// Tick sets the TickFunc event processing callback.
func (b *Builder) Tick(tick TickFunc) *Builder {
	b.tick = tick
	return b
}

// WithoutTick assigns a default no-op TickFunc callback returning an empty result.
func (b *Builder) WithoutTick() *Builder {
	b.tick = withoutTick
	return b
}

// View sets the ViewFunc rendering callback.
func (b *Builder) View(view ViewFunc) *Builder {
	b.view = view
	return b
}

func (b *Builder) makeTags() set.Set[string] {
	tags := set.New[string](5)

	if b.name == "" {
		tags.Add(ErrorMissingName)
	}

	if b.boot == nil {
		tags.Add(ErrorMissingBoot)
	}

	if b.keys == nil {
		tags.Add(ErrorMissingKeys)
	}

	if b.tick == nil {
		tags.Add(ErrorMissingTick)
	}

	if b.view == nil {
		tags.Add(ErrorMissingView)
	}

	return tags
}

func (b *Builder) makeID() string {
	if b.clock == nil {
		return b.name + "_0"
	}
	return b.name + "_" + strconv.FormatInt(b.clock(), 10)
}

func (b *Builder) toScreen() Screen {
	return Screen{
		Boot: b.boot,
		Keys: b.keys,
		Tick: b.tick,
		View: b.view,
	}
}

// ToNode compiles all builder configurations into a complete Node instance.
func (b *Builder) ToNode() Node {
	return Node{
		id:       b.makeID(),
		Name:     b.name,
		Tags:     b.makeTags(),
		Screen:   b.toScreen(),
		Stack:    b.stack,
		children: b.children,
	}
}
