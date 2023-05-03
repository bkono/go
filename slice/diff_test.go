package slice

import (
	"testing"

	"golang.org/x/exp/slices"
)

func TestDiff(t *testing.T) {
	type args[T comparable] struct {
		slice1 []T
		slice2 []T
	}
	type testCase[T comparable] struct {
		name   string
		args   args[T]
		s1Only []T
		s2Only []T
		common []T
	}
	tests := []testCase[string]{
		{
			name: "empty slices",
		},
		{
			name: "one empty slice",
			args: args[string]{
				slice1: []string{"a", "b", "c"},
			},
			s1Only: []string{"a", "b", "c"},
		},
		{
			name: "no differences",
			args: args[string]{
				slice1: []string{"a", "b", "c"},
				slice2: []string{"a", "b", "c"},
			},
			common: []string{"a", "b", "c"},
		},
		{
			name: "one difference",
			args: args[string]{
				slice1: []string{"a", "b", "c"},
				slice2: []string{"a", "b", "d"},
			},
			s1Only: []string{"c"},
			s2Only: []string{"d"},
			common: []string{"a", "b"},
		},
		{
			name: "multiple differences",
			args: args[string]{
				slice1: []string{"a", "b", "c"},
				slice2: []string{"a", "d", "e"},
			},
			s1Only: []string{"b", "c"},
			s2Only: []string{"d", "e"},
			common: []string{"a"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s1Only, s2Only, common := Diff(tt.args.slice1, tt.args.slice2)
			slices.Sort(s1Only)
			slices.Sort(s2Only)
			slices.Sort(common)
			if !slices.Equal(s1Only, tt.s1Only) {
				t.Errorf("Diff() s1Only = %v, want %v", s1Only, tt.s1Only)
			}
			if !slices.Equal(s2Only, tt.s2Only) {
				t.Errorf("Diff() s2Only = %v, want %v", s2Only, tt.s2Only)
			}
			if !slices.Equal(common, tt.common) {
				t.Errorf("Diff() common = %v, want %v", common, tt.s2Only)
			}
		})
	}
}
