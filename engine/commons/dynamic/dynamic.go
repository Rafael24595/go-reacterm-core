package dynamic

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/Rafael24595/go-reacterm-core/engine/app/hash"
)

type Value struct {
	item any
}

func From(item any) Value {
	return Value{
		item: item,
	}
}

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

func (a Value) BoolOr(def bool) bool {
	if v, ok := a.Bool(); ok {
		return v
	}
	return def
}

func (a Value) String() (string, bool) {
	switch v := a.item.(type) {
	case nil:
		return "", false
	case string:
		return v, true
	case bool:
		return strconv.FormatBool(v), false
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", v), false
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v), false
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32), false
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), false
	}
	return fmt.Sprintf("%v", a.item), false
}

func (a Value) StringOr(def string) string {
	if v, ok := a.String(); ok {
		return v
	}
	return def
}

func (a Value) Int() (int, bool) {
	if v, ok := a.Int64(); ok {
		return int(v), true
	}
	return 0, false
}

func (a Value) IntOr(def int) int {
	if v, ok := a.Int64(); ok {
		return int(v)
	}
	return def
}

func (a Value) Int8() (int8, bool) {
	if v, ok := a.Int64(); ok {
		return int8(v), true
	}
	return 0, false
}

func (a Value) Int8Or(def int8) int8 {
	if v, ok := a.Int8(); ok {
		return v
	}
	return def
}

func (a Value) Int16() (int16, bool) {
	if v, ok := a.Int64(); ok {
		return int16(v), true
	}
	return 0, false
}

func (a Value) Int16Or(def int16) int16 {
	if v, ok := a.Int16(); ok {
		return v
	}
	return def
}

func (a Value) Int32() (int32, bool) {
	if v, ok := a.Int64(); ok {
		return int32(v), true
	}
	return 0, false
}

func (a Value) Int32Or(def int32) int32 {
	if v, ok := a.Int32(); ok {
		return v
	}
	return def
}

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
		return int64(v), true
	case uint8:
		return int64(v), true
	case uint16:
		return int64(v), true
	case uint32:
		return int64(v), true
	case uint64:
		return int64(v), true
	case float32:
		return int64(v), true
	case float64:
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

func (a Value) Int64Or(def int64) int64 {
	if v, ok := a.Int64(); ok {
		return v
	}
	return def
}

func (a Value) Uint() (uint, bool) {
	if v, ok := a.Int64(); ok && v >= 0 {
		return uint(v), true
	}
	return 0, false
}

func (a Value) UintOr(def uint) uint {
	if v, ok := a.Uint(); ok {
		return v
	}
	return def
}

func (a Value) Uint8() (uint8, bool) {
	if v, ok := a.Int64(); ok && v >= 0 {
		return uint8(v), true
	}
	return 0, false
}

func (a Value) Uint8Or(def uint8) uint8 {
	if v, ok := a.Uint8(); ok {
		return v
	}
	return def
}

func (a Value) Uint16() (uint16, bool) {
	if v, ok := a.Int64(); ok && v >= 0 {
		return uint16(v), true
	}
	return 0, false
}

func (a Value) Uint16Or(def uint16) uint16 {
	if v, ok := a.Uint16(); ok {
		return v
	}
	return def
}

func (a Value) Uint32() (uint32, bool) {
	if v, ok := a.Int64(); ok && v >= 0 {
		return uint32(v), true
	}
	return 0, false
}

func (a Value) Uint32Or(def uint32) uint32 {
	if v, ok := a.Uint32(); ok {
		return v
	}
	return def
}

func (a Value) Uint64() (uint64, bool) {
	if v, ok := a.Int64(); ok && v >= 0 {
		return uint64(v), true
	}
	return 0, false
}

func (a Value) Uint64Or(def uint64) uint64 {
	if v, ok := a.Uint64(); ok {
		return v
	}
	return def
}

func (a Value) Float32() (float32, bool) {
	if v, ok := a.Float64(); ok {
		return float32(v), true
	}
	return 0, false
}

func (a Value) Float32Or(def float32) float32 {
	if v, ok := a.Float64(); ok {
		return float32(v)
	}
	return def
}

func (a Value) Float64() (float64, bool) {
	switch v := a.item.(type) {
	case int:
		return float64(v), true
	case int8:
		return float64(v), true
	case int16:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case uint:
		return float64(v), true
	case uint8:
		return float64(v), true
	case uint16:
		return float64(v), true
	case uint32:
		return float64(v), true
	case uint64:
		return float64(v), true
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

func (a Value) Float64Or(def float64) float64 {
	if v, ok := a.Float64(); ok {
		return v
	}
	return def
}

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

func Map[T any](a Value) (T, bool) {
	val, ok := a.item.(T)
	return val, ok
}

func MapOr[T any](a Value, def T) T {
	if v, ok := Map[T](a); ok {
		return v
	}
	return def
}

func Parse[T any](a Value, parse func(string) (T, error)) (T, bool) {
	var zero T

	v, err := parse(a.Text())
	if err != nil {
		return zero, false
	}

	return v, true
}

func Parsed[T any](a Value, parse func(string) (T, error), def T) T {
	if v, ok := Parse(a, parse); ok {
		return v
	}
	return def
}
