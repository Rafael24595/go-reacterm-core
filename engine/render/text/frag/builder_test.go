package frag

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
)

func TestNewBuilder(t *testing.T) {
	b := NewBuilder()

	assert.Equal(t, "", b.Text)
	assert.Equal(t, atom.None, b.Atom)
	assert.DeepEqual(t, spec.Empty(), b.Spec)
}

func TestBuilder_AddText(t *testing.T) {
	b := NewBuilder()

	b.AddText("hello")
	b.AddText(" world")

	assert.Equal(t, "hello world", b.Text)
}

func TestBuilder_AddRunes(t *testing.T) {
	b := NewBuilder()

	b.AddRunes([]rune("hello"))

	assert.Equal(t, "hello", b.Text)
}

func TestBuilder_WithMeta(t *testing.T) {
	src := New("hello", atom.Bold, spec.JustifyCenter(2))

	b := NewBuilder().
		AddText("world").
		WithMeta(&src)

	assert.Equal(t, "world", b.Text)
	assert.Equal(t, src.Atom(), b.Atom)
	assert.DeepEqual(t, src.Spec(), b.Spec)
}

func TestBuilder_WithFrag(t *testing.T) {
	src := New("golang", atom.Bold, spec.JustifyCenter(2))

	got := NewBuilder().
		AddText("hello").
		WithFrag(src).
		Frag()

	assert.Equal(t, "hellogolang", got.Text())
	assert.Equal(t, src.Atom(), got.Atom())
	assert.DeepEqual(t, src.Spec(), got.Spec())
}

func TestBuilder_AddAtom(t *testing.T) {
	b := NewBuilder()

	b.AddAtom(atom.Bold)
	b.AddAtom(atom.Focus)

	assert.Equal(t, atom.Merge(atom.Bold, atom.Focus), b.Atom)
}

func TestBuilder_AddSpec(t *testing.T) {
	b := NewBuilder()

	b.AddSpec(spec.JustifyCenter(1))
	b.AddSpec(spec.JustifyCenter(2))

	assert.DeepEqual(
		t, spec.Merge(spec.JustifyCenter(1), spec.JustifyCenter(2)), b.Spec,
	)
}

func TestBuilder_Frag(t *testing.T) {
	got := NewBuilder().
		AddText("hello").
		AddAtom(atom.Bold).
		AddSpec(spec.JustifyCenter(2)).
		Frag()

	assert.Equal(t, "hello", got.Text())
	assert.Equal(t, atom.Bold, got.Atom())
	assert.DeepEqual(t, spec.JustifyCenter(2), got.Spec())
}

func TestBuilder_Fluent(t *testing.T) {
	got := NewBuilder().
		AddText("hello").
		AddRunes([]rune(" world")).
		AddAtom(atom.Bold).
		AddAtom(atom.Focus).
		AddSpec(spec.JustifyCenter(1)).
		AddSpec(spec.JustifyCenter(2)).
		Frag()

	assert.Equal(t, "hello world", got.Text())
	assert.Equal(
		t,
		atom.Merge(atom.Bold, atom.Focus),
		got.Atom(),
	)
	assert.DeepEqual(
		t,
		spec.Merge(
			spec.JustifyCenter(1),
			spec.JustifyCenter(2),
		),
		got.Spec(),
	)
}

func TestBuilder_Frag_Hash(t *testing.T) {
	expected := New(
		"hello",
		atom.Bold,
		spec.JustifyCenter(2),
	)

	got := NewBuilder().
		AddText("hello").
		AddAtom(atom.Bold).
		AddSpec(spec.JustifyCenter(2)).
		Frag()

	assert.Equal(t, expected.hash, got.hash)
}
