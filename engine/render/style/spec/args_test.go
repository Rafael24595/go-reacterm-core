package spec

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/dynamic"
)

func TestArgsLazyInit(t *testing.T) {
	var a args

	assert.Nil(t, a.items)

	a.lazyInit()

	assert.NotNil(t, a.items)
	assert.Empty(t, a.items)
}

func TestArgsSetGet(t *testing.T) {
	var a args

	a.Set(KeyFillSize, dynamic.From(10))

	assert.Equal(
		t,
		10,
		a.Get(KeyFillSize).IntOr(0),
	)
}

func TestArgsGetMissing(t *testing.T) {
	var a args

	value := a.Get(KeyFillSize)

	assert.Equal(t, dynamic.Value{}, value)
}

func TestArgsTryGetExists(t *testing.T) {
	var a args

	a.Set(KeyFillSize, dynamic.From(15))

	v, ok := a.TryGet(KeyFillSize)

	assert.True(t, ok)
	assert.Equal(t, 15, v.IntOr(0))
}

func TestArgsTryGetMissing(t *testing.T) {
	var a args

	_, ok := a.TryGet(KeyFillSize)

	assert.False(t, ok)
}

func TestArgsDeleteExists(t *testing.T) {
	var a args

	a.Set(KeyFillSize, dynamic.From(20))

	v, ok := a.Delete(KeyFillSize)

	assert.True(t, ok)
	assert.Equal(t, 20, v.IntOr(0))

	_, ok = a.TryGet(KeyFillSize)
	assert.False(t, ok)
}

func TestArgsDeleteMissing(t *testing.T) {
	var a args

	_, ok := a.Delete(KeyFillSize)

	assert.False(t, ok)
}

func TestArgsCopy(t *testing.T) {
	src := args{}

	src.Set(KeyFillSize, dynamic.From(100))
	src.Set(KeyJustifyLeftText, dynamic.From("."))

	var dst args
	dst.Copy(src)

	assert.Equal(
		t,
		src.Get(KeyFillSize).IntOr(0),
		dst.Get(KeyFillSize).IntOr(0),
	)

	assert.Equal(
		t,
		src.Get(KeyJustifyLeftText).StringOr(""),
		dst.Get(KeyJustifyLeftText).StringOr(""),
	)
}

func TestArgsCopyOverwrite(t *testing.T) {
	src := args{}
	src.Set(KeyFillSize, dynamic.From(100))

	dst := args{}
	dst.Set(KeyFillSize, dynamic.From(50))

	dst.Copy(src)

	assert.Equal(
		t,
		100,
		dst.Get(KeyFillSize).IntOr(0),
	)
}

func TestArgsClone(t *testing.T) {
	src := args{}

	src.Set(KeyFillSize, dynamic.From(10))
	src.Set(KeyJustifyLeftText, dynamic.From("."))

	clone := src.Clone()

	assert.NotSame(t, src.items, clone.items)

	assert.Equal(
		t,
		src.Get(KeyFillSize).IntOr(0),
		clone.Get(KeyFillSize).IntOr(0),
	)

	assert.Equal(
		t,
		src.Get(KeyJustifyLeftText).IntOr(0),
		clone.Get(KeyJustifyLeftText).IntOr(0),
	)
}

func TestArgsCloneIsIndependent(t *testing.T) {
	src := args{}

	src.Set(KeyFillSize, dynamic.From(10))

	clone := src.Clone()

	clone.Set(KeyFillSize, dynamic.From(20))

	assert.Equal(
		t,
		10,
		src.Get(KeyFillSize).IntOr(0),
	)

	assert.Equal(
		t,
		20,
		clone.Get(KeyFillSize).IntOr(0),
	)
}

func BenchmarkArgsCopy(b *testing.B) {
	src := args{}

	src.Set(KeyFillSize, dynamic.From(80))
	src.Set(KeyJustifyLeftText, dynamic.From("."))
	src.Set(KeyJustifyLeftSize, dynamic.From(40))

	b.ReportAllocs()

	for b.Loop() {
		var dst args
		dst.Copy(src)
	}
}

func BenchmarkArgsClone(b *testing.B) {
	src := args{}

	src.Set(KeyFillSize, dynamic.From(80))
	src.Set(KeyJustifyLeftText, dynamic.From("."))
	src.Set(KeyJustifyLeftSize, dynamic.From(40))

	b.ReportAllocs()

	for b.Loop() {
		_ = src.Clone()
	}
}
