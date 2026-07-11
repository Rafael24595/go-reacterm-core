package dynamic

import (
	"encoding/json"
	"strconv"
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/hash"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"
	"github.com/Rafael24595/go-reacterm-core/test"
)

func TestArgumentNumericConversions(t *testing.T) {
	tests := []struct {
		name string
		from any
		want int64
	}{
		{"int to int64", int(1), 1},
		{"int8 to int64", int8(1), 1},
		{"int16 to int64", int16(1), 1},
		{"int32 to int64", int32(1), 1},
		{"int64 to int64", int64(1), 1},

		{"uint to int64", uint(1), 1},
		{"uint8 to int64", uint8(1), 1},
		{"uint16 to int64", uint16(1), 1},
		{"uint32 to int64", uint32(1), 1},
		{"uint64 to int64", uint64(1), 1},

		{"float32 to int64", float32(1.0), 1},
		{"float64 to int64", float64(1.0), 1},

		{"bool true to int64", true, 1},
		{"bool false to int64", false, 0},

		{"string '123' to int64", "123", 123},
		{"string '0' to int64", "0", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := From(tt.from).Int64()
			assert.True(t, ok)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestArgumentUnsignedConversions(t *testing.T) {
	tests := []struct {
		name string
		from any
		want uint64
	}{
		{"int to uint64", 42, 42},
		{"uint to uint64", uint(42), 42},
		{"uint8 to uint64", uint8(42), 42},
		{"uint16 to uint64", uint16(42), 42},
		{"uint32 to uint64", uint32(42), 42},
		{"uint64 to uint64", uint64(42), 42},

		{"bool true to uint64", true, 1},
		{"bool false to uint64", false, 0},

		{"string '42' to uint64", "42", 42},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := From(tt.from).Uint64()
			assert.True(t, ok)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestArgumentFloatConversions(t *testing.T) {
	tests := []struct {
		name string
		from any
		want float64
	}{
		{"int to float64", 42, 42},
		{"int64 to float64", int64(42), 42},

		{"float32 to float64", float32(42.0), 42},
		{"float64 to float64", float64(42.0), 42},

		{"bool true to float64", true, 1},
		{"bool false to float64", false, 0},

		{"string '42.5' to float64", "42.5", 42.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := From(tt.from).Float64()
			assert.True(t, ok)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestArgumentBoolConversions(t *testing.T) {
	tests := []struct {
		name string
		from any
		want bool
	}{
		{"bool true", true, true},
		{"bool false", false, false},

		{"string true", "true", true},
		{"string false", "false", false},

		{"int zero", 0, false},
		{"int nonzero", 42, true},

		{"float zero", float64(0), false},
		{"float nonzero", float64(3.14), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := From(tt.from).Bool()
			assert.True(t, ok)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestArgumentStringConversions(t *testing.T) {
	tests := []struct {
		name string
		from any
		want string
	}{
		{"nil", nil, ""},

		{"string", "hello", "hello"},

		{"int", 42, "42"},

		{"float32", float32(3.14), "3.14"},
		{"float64", float64(3.14), "3.14"},

		{"bool true", true, "true"},
		{"bool false", false, "false"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := From(tt.from).Text()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestArgumentDefaults(t *testing.T) {
	tests := []struct {
		name string
		from any
		def  any
		want any
		fn   func(Value) any
	}{
		{"Intd valid int", 42, 100, 42, func(a Value) any { return a.IntOr(100) }},
		{"Intd invalid string", "abc", 100, 100, func(a Value) any { return a.IntOr(100) }},

		{"Int64d valid int64", int64(99), 100, int64(99), func(a Value) any { return a.Int64Or(100) }},
		{"Int64d invalid string", "xyz", 100, int64(100), func(a Value) any { return a.Int64Or(100) }},

		{"Uint64d valid uint64", uint64(77), 100, uint64(77), func(a Value) any { return a.Uint64Or(100) }},
		{"Uint64d invalid string", "xyz", 100, uint64(100), func(a Value) any { return a.Uint64Or(100) }},

		{"Float32d valid float", float32(3.14), 2.71, float32(3.14), func(a Value) any { return a.Float32Or(2.71) }},
		{"Float32d invalid string", "abc", 2.71, float32(2.71), func(a Value) any { return a.Float32Or(2.71) }},

		{"Float64d valid float", float64(1.618), 3.14, float64(1.618), func(a Value) any { return a.Float64Or(3.14) }},
		{"Float64d invalid string", "abc", 3.14, float64(3.14), func(a Value) any { return a.Float64Or(3.14) }},

		{"Boold valid true", true, false, true, func(a Value) any { return a.BoolOr(false) }},
		{"Boold valid false", false, true, false, func(a Value) any { return a.BoolOr(true) }},
		{"Boold invalid string", "notbool", true, true, func(a Value) any { return a.BoolOr(true) }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arg := From(tt.from)
			got := tt.fn(arg)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMap(t *testing.T) {
	tests := []struct {
		name     string
		from     any
		wantData any
		wantOk   bool
	}{
		{"String matching", "hola", "hola", true},
		{"Int matching", 42, 42, true},
		{"[]Rune matching", []rune{'a'}, []rune{'a'}, true},
		{"Type mismatch", "100", 0, false},
		{"Nil case", nil, nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arg := From(tt.from)

			switch v := tt.wantData.(type) {
			case string:
				got, ok := Map[string](arg)

				assert.Equal(t, tt.wantOk, ok)
				assert.Equal(t, got, v)
			case int:
				got, ok := Map[int](arg)

				assert.Equal(t, tt.wantOk, ok)
				assert.Equal(t, got, v)

			case []rune:
				got, ok := Map[[]rune](arg)

				assert.Equal(t, tt.wantOk, ok)
				assert.DeepEqual(t, got, v)
			}
		})
	}
}

func TestParse_IntSuccess(t *testing.T) {
	arg := From("123")

	got, ok := Parse(arg, strconv.Atoi)

	assert.True(t, ok)
	assert.Equal(t, 123, got)
}

func TestParse_IntFailure(t *testing.T) {
	arg := From("abc")

	_, ok := Parse(arg, strconv.Atoi)

	assert.False(t, ok)
}

func TestParse_Json(t *testing.T) {
	type lang struct {
		Lang string `json:"lang"`
	}

	arg := From(`{ "lang": "golang" }`)

	parseJson := func(s string) (lang, error) {
		var l lang
		err := json.Unmarshal([]byte(s), &l)
		return l, err
	}

	got, ok := Parse(arg, parseJson)

	assert.True(t, ok)
	assert.Equal(t, "golang", got.Lang)
}

func TestArgumentHash_Types(t *testing.T) {
	tests := []struct {
		name string
		from any
	}{
		{"nil", nil},

		{"bool true", true},
		{"bool false", false},

		{"string", "hello"},

		{"int", int(42)},
		{"int8", int8(42)},
		{"int16", int16(42)},
		{"int32", int32(42)},
		{"int64", int64(42)},

		{"uint", uint(42)},
		{"uint8", uint8(42)},
		{"uint16", uint16(42)},
		{"uint32", uint32(42)},
		{"uint64", uint64(42)},

		{"float32", float32(3.14)},
		{"float64", float64(3.14)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h1 := From(tt.from).Hash(hash.New()).Sum64()
			h2 := From(tt.from).Hash(hash.New()).Sum64()

			assert.GreaterThan(t, 0, h1)
			assert.Equal(t, h1, h2)
		})
	}
}

func TestArgumentHash_DifferentTypes(t *testing.T) {
	cases := []Value{
		From(int(1)),
		From(int64(1)),
		From(uint(1)),
		From(float64(1)),
		From("1"),
		From(true),
	}

	seen := set.New[uint64]()

	for _, c := range cases {
		h := c.Hash(hash.New()).Sum64()
		assert.False(t, seen.Has(h))
		seen.Add(h)
	}
}

func TestArgumentHash_Fallback(t *testing.T) {
	a := From(test.Lang{
		Name:    "Golang",
		Version: "1.25.5",
	})

	h1 := a.Hash(hash.New()).Sum64()
	h2 := a.Hash(hash.New()).Sum64()

	assert.Equal(t, h1, h2)
}

func TestArgumentHash_IsComposable(t *testing.T) {
	h1 := hash.New().
		Uint8(7)

	h1 = From("abc").Hash(h1)

	h2 := hash.New().
		Uint8(7)

	h2 = From("abc").Hash(h2)

	assert.Equal(t, h1.Sum64(), h2.Sum64())
}
