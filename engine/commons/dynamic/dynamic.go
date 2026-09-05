package dynamic

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/Rafael24595/go-reacterm-core/engine/app/hash"
)

// Value wraps an arbitrary Go item, providing safe, multi-type casting and parsing operations.
type Value struct {
	item any
}

// From constructs a Value container around the given item.
func From(item any) Value {
	return Value{
		item: item,
	}
}

// Bool attempts to convert the underlying item into a boolean.
func (a Value) Bool() (bool, bool) {
	switch v := a.item.(type) {
	case bool:
		return v, true
	case int, int8, int16, int32, int64:
		return a.Int64Or(0) != 0, true
	case uint, uint8, uint16, uint32, uint64:
		return a.Int64Or(0) != 0, true
	case float32, float64:
		return a.Float64Or(0) != 0, true
	case string:
		val, err := strconv.ParseBool(strings.ToLower(v))
		if err == nil {
			return val, true
		}
	}
	return false, false
}

// BoolOr returns the boolean value or default if conversion fails.
func (a Value) BoolOr(def bool) bool {
	if v, ok := a.Bool(); ok {
		return v
	}
	return def
}

// String retrieves the underlying string or rune slice if available.
func (a Value) String() (string, bool) {
	switch v := a.item.(type) {
	case string:
		return v, true
	case []rune:
		return string(v), true
	}
	return "", false
}

// StringOr returns the string value or default if not a string.
func (a Value) StringOr(def string) string {
	if v, ok := a.String(); ok {
		return v
	}
	return def
}

// Int attempts to retrieve an int value.
func (a Value) Int() (int, bool) {
	if v, ok := a.Int64(); ok {
		return int(v), true
	}
	return 0, false
}

// IntOr returns the int value or default.
func (a Value) IntOr(def int) int {
	if v, ok := a.Int64(); ok {
		return int(v)
	}
	return def
}

// Int8 attempts to retrieve an int8 value.
func (a Value) Int8() (int8, bool) {
	if v, ok := a.Int64(); ok {
		return int8(v), true
	}
	return 0, false
}

// Int8Or returns the int8 value or default.
func (a Value) Int8Or(def int8) int8 {
	if v, ok := a.Int8(); ok {
		return v
	}
	return def
}

// Int16 attempts to retrieve an int16 value.
func (a Value) Int16() (int16, bool) {
	if v, ok := a.Int64(); ok {
		return int16(v), true
	}
	return 0, false
}

// Int16Or returns the int16 value or default.
func (a Value) Int16Or(def int16) int16 {
	if v, ok := a.Int16(); ok {
		return v
	}
	return def
}

// Int32 attempts to retrieve an int32 value.
func (a Value) Int32() (int32, bool) {
	if v, ok := a.Int64(); ok {
		return int32(v), true
	}
	return 0, false
}

// Int32Or returns the int32 value or default.
func (a Value) Int32Or(def int32) int32 {
	if v, ok := a.Int32(); ok {
		return v
	}
	return def
}

// Int64 converts numeric, boolean, or string representations into int64.
func (a Value) Int64() (int64, bool) {
	switch v := a.item.(type) {
	case int:
		return int64(v), true
	case int8:
		return int64(v), true
	case int16:
		return int64(v), true
	case int32:
		return int64(v), true
	case int64:
		return v, true
	case uint:
		if uint64(v) > math.MaxInt64 {
			return 0, false
		}
		return int64(v), true
	case uint8:
		return int64(v), true
	case uint16:
		return int64(v), true
	case uint32:
		return int64(v), true
	case uint64:
		if v > math.MaxInt64 {
			return 0, false
		}
		return int64(v), true
	case float32:
		if v < math.MinInt64 || v > math.MaxInt64 || math.IsNaN(float64(v)) {
			return 0, false
		}
		return int64(v), true
	case float64:
		if v < math.MinInt64 || v > math.MaxInt64 || math.IsNaN(v) {
			return 0, false
		}
		return int64(v), true
	case bool:
		if v {
			return 1, true
		}
		return 0, true
	case string:
		val, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			return val, true
		}
	}
	return 0, false
}

// Int64Or returns the int64 value or default.
func (a Value) Int64Or(def int64) int64 {
	if v, ok := a.Int64(); ok {
		return v
	}
	return def
}

// Uint attempts to retrieve a uint value.
func (a Value) Uint() (uint, bool) {
	if v, ok := a.Int64(); ok && v >= 0 {
		return uint(v), true
	}
	return 0, false
}

// UintOr returns the uint value or default.
func (a Value) UintOr(def uint) uint {
	if v, ok := a.Uint(); ok {
		return v
	}
	return def
}

// Uint8 attempts to retrieve a uint8 value.
func (a Value) Uint8() (uint8, bool) {
	if v, ok := a.Int64(); ok && v >= 0 {
		return uint8(v), true
	}
	return 0, false
}

// Uint8Or returns the uint8 value or default.
func (a Value) Uint8Or(def uint8) uint8 {
	if v, ok := a.Uint8(); ok {
		return v
	}
	return def
}

// Uint16 attempts to retrieve a uint16 value.
func (a Value) Uint16() (uint16, bool) {
	if v, ok := a.Int64(); ok && v >= 0 {
		return uint16(v), true
	}
	return 0, false
}

// Uint16Or returns the uint16 value or default.
func (a Value) Uint16Or(def uint16) uint16 {
	if v, ok := a.Uint16(); ok {
		return v
	}
	return def
}

// Uint32 attempts to retrieve a uint32 value.
func (a Value) Uint32() (uint32, bool) {
	if v, ok := a.Int64(); ok && v >= 0 {
		return uint32(v), true
	}
	return 0, false
}

// Uint32Or returns the uint32 value or default.
func (a Value) Uint32Or(def uint32) uint32 {
	if v, ok := a.Uint32(); ok {
		return v
	}
	return def
}

// Uint64 attempts to retrieve a uint64 value.
func (a Value) Uint64() (uint64, bool) {
	switch v := a.item.(type) {
	case uint64:
		return v, true
	case uint, uint8, uint16, uint32:
		return uint64(a.Int64Or(0)), true
	case int, int8, int16, int32, int64:
		val := a.Int64Or(-1)
		if val >= 0 {
			return uint64(val), true
		}
	case bool:
		if v {
			return 1, true
		}
		return 0, true
	case string:
		val, err := strconv.ParseUint(v, 10, 64)
		if err == nil {
			return val, true
		}
	}
	return 0, false
}

// Uint64Or returns the uint64 value or default.
func (a Value) Uint64Or(def uint64) uint64 {
	if v, ok := a.Uint64(); ok {
		return v
	}
	return def
}

// Float32 attempts to retrieve a float32 value.
func (a Value) Float32() (float32, bool) {
	if v, ok := a.Float64(); ok {
		return float32(v), true
	}
	return 0, false
}

// Float32Or returns the float32 value or default.
func (a Value) Float32Or(def float32) float32 {
	if v, ok := a.Float64(); ok {
		return float32(v)
	}
	return def
}

// Float64 converts numeric, boolean, or string representations into float64.
func (a Value) Float64() (float64, bool) {
	switch v := a.item.(type) {
	case int, int8, int16, int32, int64:
		return float64(a.Int64Or(0)), true
	case uint, uint8, uint16, uint32, uint64:
		if u, ok := a.Uint64(); ok {
			return float64(u), true
		}
	case float32:
		return float64(v), true
	case float64:
		return v, true
	case bool:
		if v {
			return 1, true
		}
		return 0, true
	case string:
		val, err := strconv.ParseFloat(v, 64)
		if err == nil {
			return val, true
		}
	}
	return 0, false
}

// Float64Or returns the float64 value or default.
func (a Value) Float64Or(def float64) float64 {
	if v, ok := a.Float64(); ok {
		return v
	}
	return def
}

// Text converts any primitive type to its string representation.
func (a Value) Text() string {
	switch v := a.item.(type) {
	case nil:
		return ""
	case string:
		return v
	case bool:
		return strconv.FormatBool(v)
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", v)
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	}
	return fmt.Sprintf("%v", a.item)
}

// Hash feeds the internal type tag and bits into a provided Hasher instance.
func (a Value) Hash(h hash.Hasher) hash.Hasher {
	switch v := a.item.(type) {
	case nil:
		return h.Uint8(Nil.Uint8()).
			Uint8(0)

	case bool:
		return h.Uint8(Bool.Uint8()).
			Bool(v)

	case string:
		return h.Uint8(String.Uint8()).
			String(v)

	case int:
		return h.Uint8(Int.Uint8()).
			Uint64(uint64(v))

	case int8:
		return h.Uint8(Int8.Uint8()).
			Uint8(uint8(v))

	case int16:
		return h.Uint8(Int16.Uint8()).
			Uint16(uint16(v))

	case int32:
		return h.Uint8(Int32.Uint8()).
			Uint32(uint32(v))

	case int64:
		return h.Uint8(Int64.Uint8()).
			Uint64(uint64(v))

	case uint:
		return h.Uint8(Uint.Uint8()).
			Uint64(uint64(v))

	case uint8:
		return h.Uint8(Uint8.Uint8()).
			Uint8(v)

	case uint16:
		return h.Uint8(Uint16.Uint8()).
			Uint16(v)

	case uint32:
		return h.Uint8(Uint32.Uint8()).
			Uint32(v)

	case uint64:
		return h.Uint8(Uint64.Uint8()).
			Uint64(v)

	case float32:
		return h.Uint8(Float32.Uint8()).
			Uint32(math.Float32bits(v))

	case float64:
		return h.Uint8(Float64.Uint8()).
			Uint64(math.Float64bits(v))
	}

	return h.Uint8(Fallback.Uint8()).
		String(a.Text())
}

// Map attempts a direct generic type-assertion to type T.
func Map[T any](a Value) (T, bool) {
	val, ok := a.item.(T)
	return val, ok
}

// MapOr attempts type-assertion to T or returns the default value.
func MapOr[T any](a Value, def T) T {
	if v, ok := Map[T](a); ok {
		return v
	}
	return def
}

// Parse applies a custom parsing function over the text representation of the Value.
func Parse[T any](a Value, parse func(string) (T, error)) (T, bool) {
	var zero T

	v, err := parse(a.Text())
	if err != nil {
		return zero, false
	}

	return v, true
}

// Parsed applies a custom parser or returns the default value if parsing fails.
func Parsed[T any](a Value, parse func(string) (T, error), def T) T {
	if v, ok := Parse(a, parse); ok {
		return v
	}
	return def
}
