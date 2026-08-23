package styler

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-reacterm-core/engine/format"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
)

func TestSpecRegistryCoverage(t *testing.T) {
	for r := range spec.Registry() {
		var found bool

		for _, s := range specs {
			if s.Kind == r.Kind() {
				found = true
				break
			}
		}

		assert.True(t, found, "spec %s not found", r.Name)
	}
}

func TestSpecsAreUnique(t *testing.T) {
	cache := make(map[spec.Kind]bool)

	for _, v := range specs {
		_, ok := cache[v.Kind]
		assert.False(t, ok)

		cache[v.Kind] = true
	}
}

func TestDeduplicateSpecs(t *testing.T) {
	ruleFill := SpecRule{
		Kind: spec.KindFill,
		Fn: func(spec.Spec, winsize.Cols, format.Text) (string, bool) {
			return "rule_fill", false
		},
	}

	ruleTruncLeft := SpecRule{
		Kind: spec.KindTruncateLeft,
		Fn: func(spec.Spec, winsize.Cols, format.Text) (string, bool) {
			return "rule_trunc_left", false
		},
	}

	ruleTruncRight := SpecRule{
		Kind: spec.KindTruncateRight,
		Fn: func(spec.Spec, winsize.Cols, format.Text) (string, bool) {
			return "rule_trunc_right", false
		},
	}

	ruleMock := func(kind spec.Kind, text string) SpecRule {
		return SpecRule{
			Kind: kind,
			Fn: func(spec.Spec, winsize.Cols, format.Text) (string, bool) {
				return text, false
			},
		}
	}

	tests := []struct {
		name  string
		input []SpecRule
		want  []string
	}{
		{
			name:  "empty",
			input: nil,
			want:  []string{},
		},
		{
			name: "without duplicates",
			input: []SpecRule{
				ruleFill,
				ruleTruncLeft,
				ruleTruncRight,
			},
			want: []string{
				"rule_fill",
				"rule_trunc_left",
				"rule_trunc_right",
			},
		},
		{
			name: "replaces duplicate with last rule",
			input: []SpecRule{
				ruleFill,
				ruleTruncLeft,
				ruleMock(spec.KindFill, "rule_mock"),
			},
			want: []string{
				"rule_mock",
				"rule_trunc_left",
			},
		},
		{
			name: "preserves first occurrence position",
			input: []SpecRule{
				ruleFill,
				ruleTruncLeft,
				ruleTruncRight,
				ruleMock(spec.KindFill, "rule_mock"),
			},
			want: []string{
				"rule_mock",
				"rule_trunc_left",
				"rule_trunc_right",
			},
		},
		{
			name: "handles multiple duplicates",
			input: []SpecRule{
				ruleFill,
				ruleTruncLeft,
				ruleMock(spec.KindFill, "rule_mock_01"),
				ruleTruncRight,
				ruleMock(spec.KindTruncateLeft, "rule_mock_02"),
				ruleMock(spec.KindFill, "rule_mock_03"),
			},
			want: []string{
				"rule_mock_03",
				"rule_mock_02",
				"rule_trunc_right",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := deduplicateSpecs(tt.input)

			assert.Size(t, len(tt.want), got)

			for i := range tt.want {
				want, _ := got[i].Fn(spec.Empty(), 0, format.Text{})
				assert.Equal(t, want, tt.want[i])
			}
		})
	}
}
