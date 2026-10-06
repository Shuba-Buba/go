package collect

import (
	"reflect"
	"testing"
)

func TestCollect(t *testing.T) {
	ch := make(chan int, 3)
	ch <- 2
	ch <- 3
	ch <- 5
	close(ch)

	got := Collect(ch)
	want := []int{2, 3, 5}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Collect = %v, want %v", got, want)
	}
}

func TestCollectEmpty(t *testing.T) {
	ch := make(chan int)
	close(ch)
	got := Collect(ch)
	if got == nil || len(got) != 0 {
		t.Fatalf("Collect(empty) = %#v, want non-nil empty slice", got)
	}
}
