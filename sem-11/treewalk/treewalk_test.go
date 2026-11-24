package treewalk

import (
	"slices"
	"testing"
)

func TestVisitLeaves(t *testing.T) {
	root := &Tree[string]{
		Value: "раздел",
		Left:  &Tree[string]{Value: "текст"},
		Right: &Tree[string]{
			Value: "приложения",
			Left:  &Tree[string]{Value: "текст"},
			Right: &Tree[string]{Right: &Tree[string]{Value: "таблица"}},
		},
	}

	var got []string

	VisitLeaves(root, func(value string) { got = append(got, value) })

	want := []string{"текст", "текст", "таблица"}
	if !slices.Equal(got, want) {
		t.Fatalf("visited %q, want %q", got, want)
	}
}

func TestEmptyAndSingleLeaf(t *testing.T) {
	VisitLeaves(nil, func(int) { t.Fatal("visited a node in an empty tree") })

	calls := 0

	VisitLeaves(&Tree[[]int]{Value: []int{1, 2}}, func(value []int) {
		calls++

		if !slices.Equal(value, []int{1, 2}) {
			t.Fatalf("leaf = %v, want [1 2]", value)
		}
	})

	if calls != 1 {
		t.Fatalf("visited single leaf %d times", calls)
	}
}
