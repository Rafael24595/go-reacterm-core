package line

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
)

func TestNewBuilder(t *testing.T) {
	b := NewBuilder()

	assert.Equal(t, 0, b.Order)
	assert.Equal(t, spec.Empty().Hash(), b.Spec.Hash())
	assert.Empty(t, b.Text)
}

func TestNewBuilderWithCapacity(t *testing.T) {
	b := NewBuilder(10)

	assert.Empty(t, b.Text)
	assert.Equal(t, 10, cap(b.Text))
}

func TestBuilderSetOrder(t *testing.T) {
	b := NewBuilder().
		SetOrder(42)

	assert.Equal(t, 42, b.Order)
}

func TestBuilderSetSpec(t *testing.T) {
	b := NewBuilder().
		AddSpec(spec.JustifyRight(5, ".")).
		SetSpec(spec.Fill(80))

	assert.Equal(t, spec.Fill(80).Hash(), b.Spec.Hash())
}

func TestBuilderAddSpec(t *testing.T) {
	b := NewBuilder().
		SetSpec(spec.Fill(80)).
		AddSpec(spec.JustifyRight(5, "."))

	expected := spec.Merge(
		spec.Fill(80),
		spec.JustifyRight(5, "."),
	)

	assert.Equal(t, expected.Hash(), b.Spec.Hash())
}

func TestBuilderPushText(t *testing.T) {
	b := NewBuilder().
		PushText("A", "B")

	assert.Size(t, 2, b.Text)
	assert.Equal(t, "A", b.Text[0].Text())
	assert.Equal(t, "B", b.Text[1].Text())
}

func TestBuilderUnshiftText(t *testing.T) {
	b := NewBuilder().
		PushText("C").
		UnshiftText("A", "B")

	assert.Equal(t, 3, len(b.Text))
	assert.Equal(t, "A", b.Text[0].Text())
	assert.Equal(t, "B", b.Text[1].Text())
	assert.Equal(t, "C", b.Text[2].Text())
}

func TestBuilderPushFrags(t *testing.T) {
	b := NewBuilder().
		PushFrags(
			frag.FromString("A"),
			frag.FromString("B"),
		)

	assert.Equal(t, 2, len(b.Text))
	assert.Equal(t, "A", b.Text[0].Text())
	assert.Equal(t, "B", b.Text[1].Text())
}

func TestBuilderUnshiftFrags(t *testing.T) {
	builder := NewBuilder().
		PushFrags(
			frag.FromString("C"),
		).
		UnshiftFrags(
			frag.FromString("A"),
			frag.FromString("B"),
		)

	assert.Equal(t, "A", builder.Text[0].Text())
	assert.Equal(t, "B", builder.Text[1].Text())
	assert.Equal(t, "C", builder.Text[2].Text())
}

func TestBuilderPushBuilder(t *testing.T) {
	a := frag.NewBuilder().
		AddText("A")

	bb := frag.NewBuilder().
		AddText("B")

	builder := NewBuilder().
		PushBuilder(a, bb)

	assert.Size(t, 2, builder.Text)
	assert.Equal(t, "A", builder.Text[0].Text())
	assert.Equal(t, "B", builder.Text[1].Text())
}

func TestBuilderUnshiftBuilder(t *testing.T) {
	builder := NewBuilder().
		PushBuilder(
			frag.NewBuilder().AddText("C"),
		).
		UnshiftBuilder(
			frag.NewBuilder().AddText("A"),
			frag.NewBuilder().AddText("B"),
		)

	assert.Equal(t, "A", builder.Text[0].Text())
	assert.Equal(t, "B", builder.Text[1].Text())
	assert.Equal(t, "C", builder.Text[2].Text())
}

func TestBuilderWithMeta(t *testing.T) {
	line := NewBuilder().
		SetOrder(12).
		SetSpec(spec.Fill(80)).
		PushText("Hello").
		Line()

	builder := NewBuilder().
		WithMeta(line)

	assert.Equal(t, 12, builder.Order)
	assert.Equal(t, line.Spec.Hash(), builder.Spec.Hash())
	assert.Empty(t, builder.Text)
}

func TestBuilderWithLine(t *testing.T) {
	line := NewBuilder().
		SetOrder(8).
		SetSpec(spec.Fill(40)).
		PushText("Hello", "World").
		Line()

	builder := NewBuilder().
		WithLine(line)

	assert.Equal(t, uint16(8), builder.Order)
	assert.Equal(t, line.Spec.Hash(), builder.Spec.Hash())

	assert.Size(t, 2, builder.Text)
	assert.Equal(t, "Hello", builder.Text[0].Text())
	assert.Equal(t, "World", builder.Text[1].Text())
}

func TestBuilderLine(t *testing.T) {
	builder := NewBuilder().
		SetOrder(5).
		SetSpec(spec.Fill(20)).
		PushText("Hello")

	line := builder.Line()

	assert.Equal(t, uint16(5), line.Order)
	assert.Equal(t, spec.Fill(20).Hash(), line.Spec.Hash())

	assert.Equal(t, 1, len(line.Text))
	assert.Equal(t, "Hello", line.Text[0].Text())
}

func TestBuilderLinePtr(t *testing.T) {
	builder := NewBuilder().
		SetOrder(10).
		PushText("Hello")

	ref := builder.LinePtr()
	val := builder.Line()

	assert.NotNil(t, ref)

	assert.Equal(t, val.Order, ref.Order)
	assert.DeepEqual(t, val.Text, ref.Text)
}

func TestBuilderLineIsImmutable(t *testing.T) {
	builder := NewBuilder().
		PushText("Hello")

	line := builder.Line()

	builder.PushText("World")

	assert.Equal(t, 1, len(line.Text))
	assert.Equal(t, "Hello", line.Text[0].Text())
}
