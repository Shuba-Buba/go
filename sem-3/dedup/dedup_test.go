package dedup

import (
	"reflect"
	"testing"
)

func TestDeduplicate(t *testing.T) {
	testCases := []struct {
		name  string
		input []string
		want  []string
	}{
		{
			name:  "duplicates",
			input: []string{"go", "map", "go", "slice", "map"},
			want:  []string{"go", "map", "slice"},
		},
		{
			name:  "order of first occurrences",
			input: []string{"b", "a", "b", "c", "a"},
			want:  []string{"b", "a", "c"},
		},
		{
			name:  "empty strings",
			input: []string{"", "", "go", ""},
			want:  []string{"", "go"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			before := append([]string(nil), tc.input...)

			got := Deduplicate(tc.input)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("Deduplicate(%v) = %v, want %v", tc.input, got, tc.want)
			}

			if !reflect.DeepEqual(tc.input, before) {
				t.Fatalf("Deduplicate changed input: got %v, want %v", tc.input, before)
			}
		})
	}
}

func TestDeduplicateNilAndEmpty(t *testing.T) {
	if got := Deduplicate(nil); got != nil {
		t.Fatalf("Deduplicate(nil) = %#v, want nil", got)
	}

	empty := []string{}

	got := Deduplicate(empty)
	if got == nil || len(got) != 0 {
		t.Fatalf("Deduplicate([]) = %#v, want a non-nil empty slice", got)
	}
}
