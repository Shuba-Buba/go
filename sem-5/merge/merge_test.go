package merge

import (
	"reflect"
	"sort"
	"testing"
	"time"
)

func send(delay time.Duration, vals ...int) <-chan int {
	ch := make(chan int)
	go func() {
		if delay > 0 {
			time.Sleep(delay)
		}
		for _, v := range vals {
			ch <- v
		}
		close(ch)
	}()
	return ch
}

func collect(in <-chan int) []int {
	out := make([]int, 0)
	for v := range in {
		out = append(out, v)
	}
	return out
}

func sorted(in []int) []int {
	out := append([]int(nil), in...)
	sort.Ints(out)
	return out
}

func TestMergeTwoSources(t *testing.T) {
	got := collect(Merge(send(0, 1, 2), send(10*time.Millisecond, 3, 4)))
	want := []int{1, 2, 3, 4}
	if !reflect.DeepEqual(sorted(got), want) {
		t.Fatalf("Merge = %v, want multiset %v", got, want)
	}
}

func TestMergeSingleAndEmptyInput(t *testing.T) {
	got := collect(Merge(send(0, 7, 8)))
	if !reflect.DeepEqual(got, []int{7, 8}) {
		t.Fatalf("Merge(single) = %v", got)
	}

	closed := make(chan int)
	close(closed)
	got = collect(Merge(closed))
	if got == nil || len(got) != 0 {
		t.Fatalf("Merge(closed empty) = %v, want empty slice", got)
	}
}

func TestMergeNoInputs(t *testing.T) {
	got := collect(Merge())
	if got == nil || len(got) != 0 {
		t.Fatalf("Merge() = %v, want empty slice", got)
	}
}

func TestMergeDoesNotCloseInputs(t *testing.T) {
	ch := make(chan int, 1)
	ch <- 42
	_ = Merge(ch)
	select {
	case _, ok := <-ch:
		if !ok {
			t.Fatal("Merge closed input channel")
		}
	default:
		t.Fatal("expected to read remaining value from input")
	}
}
