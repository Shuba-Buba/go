package timeline

import (
	"errors"
	"reflect"
	"testing"
)

func TestBuild(t *testing.T) {
	changes := []Change{
		{Day: 3, Delta: -1},
		{Day: 1, Delta: 4},
		{Day: 3, Delta: 2},
		{Day: 2, Delta: 5},
		{Day: 2, Delta: -5},
	}
	want := []Snapshot{
		{Day: 1, Value: 6},
		{Day: 3, Value: 7},
	}

	before := append([]Change(nil), changes...)

	got, err := Build(2, changes)
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Build returned %v, want %v", got, want)
	}

	if !reflect.DeepEqual(changes, before) {
		t.Fatalf("Build changed input: got %v, want %v", changes, before)
	}
}

func TestBuildNegativeValue(t *testing.T) {
	changes := []Change{
		{Day: 1, Delta: 2},
		{Day: 2, Delta: -4},
	}

	got, err := Build(1, changes)
	if !errors.Is(err, ErrNegativeValue) {
		t.Fatalf("Build error = %v, want ErrNegativeValue", err)
	}

	if got != nil {
		t.Fatalf("Build returned partial result %v, want nil", got)
	}

	got, err = Build(-1, nil)
	if !errors.Is(err, ErrNegativeValue) || got != nil {
		t.Fatalf("Build(-1, nil) = (%v, %v), want (nil, ErrNegativeValue)", got, err)
	}
}

func TestBuildEmpty(t *testing.T) {
	got, err := Build(10, nil)
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	if got == nil || len(got) != 0 {
		t.Fatalf("Build(10, nil) = %#v, want a non-nil empty slice", got)
	}
}

func TestBuildAppliesNetChangePerDay(t *testing.T) {
	changes := []Change{
		{Day: 1, Delta: -3},
		{Day: 1, Delta: 3},
	}

	got, err := Build(0, changes)
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	if got == nil || len(got) != 0 {
		t.Fatalf("Build returned %#v, want a non-nil empty slice", got)
	}
}
