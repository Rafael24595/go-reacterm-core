package math

import (
	"cmp"
)

// Number represents all generic integer primitives.
type Number interface {
	Signed | Unsigned
}

// Signed constraints underlying Go signed integer types.
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Unsigned constraints underlying Go unsigned integer types.
type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// Abs returns the absolute value of val.
func Abs[T Number](val T) T {
	if val < 0 {
		return -val
	}
	return val
}

// MinNotZero returns the non-zero minimum between x and y.
// If one argument is zero, the other argument is returned.
func MinNotZero[T Number](x, y T) T {
	if x == 0 {
		return y
	}

	if y == 0 {
		return x
	}

	return min(x, y)
}

// Clamp restricts val between lower and upper limits.
func Clamp[T cmp.Ordered](val, lower, upper T) T {
	return max(lower, min(val, upper))
}

// SubClampZero subtracts b from a, returning 0 if the result would be negative.
func SubClampZero[T Number](a, b T) T {
	if a < b {
		return 0
	}
	return a - b
}

// SubClampZeroAs subtracts b from a (clamped to zero) and casts the result to type K.
func SubClampZeroAs[T Number, K Number](a, b T) K {
	return K(SubClampZero(a, b))
}

// Digits calculates the number of decimal digits in val (base 10).
func Digits[T Number](val T) uint32 {
	if val == 0 {
		return 1
	}

	fix := Abs(val) / 10
	count := uint32(1)

	for fix > 0 {
		fix = fix / 10
		count++
	}

	return count
}

// Sum accumulates all items in the slice and returns total.
func Sum[T Number](items []T) T {
	total := T(0)
	for _, item := range items {
		total += item
	}
	return total
}

// MaxMap finds the key-value pair with the highest value in a map.
// Returns false if the map is empty.
func MaxMap[K comparable, V cmp.Ordered](m map[K]V) (K, V, bool) {
	if len(m) == 0 {
		var zeroK K
		var zeroV V
		return zeroK, zeroV, false
	}

	var maxK K
	var maxV V
	var init bool

	for k, v := range m {
		if v < maxV && init {
			continue
		}

		maxK = k
		maxV = v

		init = true
	}

	return maxK, maxV, true
}
