package format

import (
	"strings"

	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
)

// Index defines a functional generator for list item indicators.
// - cursor: the index of the currently active/selected item.
// - index: the current item position being rendered.
type Index = func(cursor, index int) string

// IndexProvider defines a factory that receives column constraints and builds an Index generator.
type IndexProvider = func(digits winsize.Cols) Index

const DefaultSuffix = ".- "

var (
	// GreaterIndex emits ">" for idle items and "-" when the item is active (cursor == index).
	GreaterIndex = CustomIndex("-", ">")
	// HyphenIndex emits "-" for idle items and ">" as a selection pointer when the item is active (cursor == index).
	HyphenIndex = CustomIndex(">", "-")
)

// NumericListIndex constructs an IndexProvider formatted for standard lists (e.g. "1.- ").
func NumericListIndex() IndexProvider {
	return NumericIndex(DefaultSuffix)
}

// NumericIndex returns an IndexProvider that formats items sequentially as numbers (1-based).
func NumericIndex(suffix ...string) IndexProvider {
	filler := marker.DefaultPaddingText
	joinSufix := strings.Join(suffix, "")

	return func(digits winsize.Cols) Index {
		return func(_, index int) string {
			text := TextFromAny(index + 1)
			return JustifyLeft(digits, text, filler) + joinSufix
		}
	}
}

// AlphabeticListIndex constructs an IndexProvider formatted for standard alphabetic lists (e.g. "a.- ").
func AlphabeticListIndex() IndexProvider {
	return AlphabeticIndex(DefaultSuffix)
}

// AlphabeticIndex returns an IndexProvider that formats items sequentially using letters (a, b, c...).
func AlphabeticIndex(suffix ...string) IndexProvider {
	filler := marker.DefaultPaddingText
	joinSufix := strings.Join(suffix, "")

	return func(digits winsize.Cols) Index {
		return func(_, index int) string {
			alpha := NumberToAlpha(index + 1)
			text := TextFromAny(alpha)
			return JustifyLeft(digits, text, filler) + joinSufix
		}
	}
}

// CustomIndex returns an IndexProvider that ignores column width constraints
// and emits activeSymbol if the item is selected or idleSymbol otherwise.
func CustomIndex(activeSymbol string, idleSymbol string) IndexProvider {
	return func(_ winsize.Cols) Index {
		return func(cursor, index int) string {
			if cursor == index {
				return activeSymbol
			}
			return idleSymbol
		}
	}
}
