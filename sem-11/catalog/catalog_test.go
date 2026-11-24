package catalog

import (
	"slices"
	"testing"
)

type label string

type exhibit struct {
	key  label
	tags []string
}

func (e exhibit) ID() label { return e.key }

func TestIndexByID(t *testing.T) {
	records := []exhibit{{"a", []string{"vase"}}, {"b", []string{"coin"}}}

	index, ok := IndexByID(records)
	if !ok || len(index) != len(records) {
		t.Fatalf("index = %v, ok = %t", index, ok)
	}

	for _, record := range records {
		got, found := index[record.key]
		if !found || got.key != record.key || !slices.Equal(got.tags, record.tags) {
			t.Fatalf("index[%q] = %v, want %v", record.key, got, record)
		}
	}
}

type location struct{ room, number int }

func (l location) ID() location { return l }

func TestDuplicateAndEmpty(t *testing.T) {
	index, ok := IndexByID([]location{{}, {1, 1}, {}})
	if ok || index != nil {
		t.Fatalf("duplicate ID: index = %v, ok = %t, want nil, false", index, ok)
	}

	index, ok = IndexByID([]location(nil))
	if !ok || index == nil || len(index) != 0 {
		t.Fatalf("empty input: index = %v, ok = %t, want empty writable map, true", index, ok)
	}

	key := location{1, 2}
	index[key] = key
}
