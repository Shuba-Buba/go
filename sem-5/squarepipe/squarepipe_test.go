package squarepipe

import (
	"reflect"
	"testing"
)

func send(vals ...int) <-chan int {
	ch := make(chan int)
	go func() {
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

func TestSquare(t *testing.T) {
	got := collect(Square(send(0, 1, 2, 3, 4)))
	want := []int{0, 1, 4, 9, 16}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Square = %v, want %v", got, want)
	}
}

func TestSquareEmpty(t *testing.T) {
	in := make(chan int)
	close(in)
	got := collect(Square(in))
	if got == nil || len(got) != 0 {
		t.Fatalf("Square(empty) = %v, want empty slice", got)
	}
}
