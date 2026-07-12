package spec

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/dynamic"
)

func argsFrom(key ArgKey, value dynamic.Value) args {
	instance := args{}
	instance.Set(key, value)
	return instance
}

func assertSpec(
	t *testing.T,
	got Spec,
	wantKind Kind,
	want map[ArgKey]dynamic.Value,
) {
	t.Helper()

	assert.Equal(t, wantKind, got.Kind())
	assertArgs(t, got.args.items, want)
}

func assertArgs(
	t *testing.T,
	got map[ArgKey]dynamic.Value,
	want map[ArgKey]dynamic.Value,
) {
	t.Helper()

	assert.Equal(t, len(want), len(got))

	for key, expected := range want {
		actual, ok := got[key]
		assert.True(t, ok)
		assert.Equal(t, expected, actual)
	}
}

func TestNewBuilder_Empty(t *testing.T) {
	b := newBuilder()

	assert.Equal(t, KindNone, b.kind)
	assert.Size(t, 0, b.args.items)
}

func TestBuilderApply_AddsKinds(t *testing.T) {
	left := New(
		KindFill,
		argsFrom(
			KeyFillSize, dynamic.From(10),
		),
	)

	right := New(
		KindJustifyLeft,
		argsFrom(
			KeyJustifyLeftText, dynamic.From("..."),
		),
	)

	b := newBuilder().
		apply(left, right)

	assert.True(t, b.has(KindFill))
	assert.True(t, b.has(KindJustifyLeft))
}

func TestBuilderApply_LastArgumentWins(t *testing.T) {
	first := New(
		KindFill,
		argsFrom(
			KeyFillSize, dynamic.From(10),
		),
	)

	second := New(
		KindFill,
		argsFrom(
			KeyFillSize, dynamic.From(20),
		),
	)

	b := newBuilder().
		apply(first, second)

	value, ok := b.args.TryGet(KeyFillSize)

	assert.True(t, ok)
	assert.Equal(t, 20, dynamic.MapOr(value, 0))
}

func TestBuilderAdd(t *testing.T) {
	b := newBuilder()

	b.add(KindFill)

	assert.True(t, b.has(KindFill))
}

func TestBuilderRemove(t *testing.T) {
	b := newBuilder()

	b.add(KindFill).
		add(KindJustifyLeft)

	b.remove(KindFill)

	assert.False(t, b.has(KindFill))
	assert.True(t, b.has(KindJustifyLeft))
}

func TestBuilderSet(t *testing.T) {
	b := newBuilder()

	b.set(KeyFillSize, dynamic.From(12))

	value, ok := b.args.TryGet(KeyFillSize)

	assert.True(t, ok)
	assert.Equal(t, 12, dynamic.MapOr(value, 0))
}

func TestBuilderDelete(t *testing.T) {
	b := newBuilder()

	b.set(KeyFillSize, dynamic.From(10))
	b.delete(KeyFillSize)

	_, ok := b.args.TryGet(KeyFillSize)

	assert.False(t, ok)
}

func TestBuilderErase_RemovesKindAndArguments(t *testing.T) {
	spec := Merge(
		New(
			KindFill,
			argsFrom(
				KeyFillSize, dynamic.From(10),
			),
		),
		New(
			KindJustifyLeft,
			argsFrom(
				KeyJustifyLeftText, dynamic.From("..."),
			),
		),
	)

	builder := newBuilder(spec)

	removed := builder.erase(KindFill)

	assert.False(t, builder.has(KindFill))
	assert.True(t, builder.has(KindJustifyLeft))

	_, ok := builder.args.TryGet(KeyFillSize)
	assert.False(t, ok)

	assert.Equal(t, KindFill, removed.Kind())

	value, ok := removed.args.TryGet(KeyFillSize)
	assert.True(t, ok)
	assert.Equal(t, 10, dynamic.MapOr(value, 0))
}

func TestBuilderErase_IgnoresMissingKinds(t *testing.T) {
	spec := New(
		KindFill,
		argsFrom(
			KeyFillSize, dynamic.From(10),
		),
	)

	builder := newBuilder(spec)

	removed := builder.erase(KindJustifyLeft)

	assert.Equal(t, KindNone, removed.Kind())
	assert.True(t, builder.has(KindFill))
}

func TestBuilderBuild_EqualsMerge(t *testing.T) {
	left := New(
		KindFill,
		argsFrom(
			KeyFillSize, dynamic.From(10),
		),
	)

	right := New(
		KindJustifyLeft,
		argsFrom(
			KeyJustifyLeftText, dynamic.From("..."),
		),
	)

	got := newBuilder().
		apply(left, right).
		build()

	want := Merge(left, right)

	assertSpec(t, want, got.kind, got.args.items)
}
