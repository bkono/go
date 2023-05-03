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
		name  string
		args  args[T]
		diff1 []T
		diff2 []T
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
			diff1: []string{"a", "b", "c"},
		},
		{
			name: "no differences",
			args: args[string]{
				slice1: []string{"a", "b", "c"},
				slice2: []string{"a", "b", "c"},
			},
		},
		{
			name: "one difference",
			args: args[string]{
				slice1: []string{"a", "b", "c"},
				slice2: []string{"a", "b", "d"},
			},
			diff1: []string{"c"},
			diff2: []string{"d"},
		},
		{
			name: "multiple differences",
			args: args[string]{
				slice1: []string{"a", "b", "c"},
				slice2: []string{"a", "d", "e"},
			},
			diff1: []string{"b", "c"},
			diff2: []string{"d", "e"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1, got2 := Diff(tt.args.slice1, tt.args.slice2)
			slices.Sort(got1)
			slices.Sort(got2)
			if !slices.Equal(got1, tt.diff1) {
				t.Errorf("Diff() got1 = %v, want %v", got1, tt.diff1)
			}
			if !slices.Equal(got2, tt.diff2) {
				t.Errorf("Diff() got2 = %v, want %v", got2, tt.diff2)
			}
		})
	}
}
