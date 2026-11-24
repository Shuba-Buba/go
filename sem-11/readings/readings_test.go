package readings

import (
	"slices"
	"testing"
)

func TestRecordHighs(t *testing.T) {
	type (
		degree int16
		series []degree
	)

	original := series{-4, -5, -1, -1, -2, 3}
	input := slices.Clone(original)
	got := RecordHighs(input)

	if _, ok := any(got).(series); !ok {
		t.Fatalf("result type = %T, want series", got)
	}

	if want := (series{-4, -1, 3}); !slices.Equal(got, want) {
		t.Fatalf("RecordHighs(%v) = %v, want %v", original, got, want)
	}

	got[0] = 99

	if !slices.Equal(input, original) {
		t.Fatalf("input changed to %v, want %v", input, original)
	}
}

func TestStringsAndEmptyInput(t *testing.T) {
	got := RecordHighs([]string{"b", "a", "d"})
	if !slices.Equal(got, []string{"b", "d"}) {
		t.Fatalf("string records = %v, want [b d]", got)
	}

	var empty []float64
	if got := RecordHighs(empty); len(got) != 0 {
		t.Fatalf("empty input returned %v", got)
	}
}
