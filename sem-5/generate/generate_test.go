package generate

import (
	"reflect"
	"testing"
)

func TestGenerate(t *testing.T) {
	got := Collect(Generate(5))
	want := []int{0, 1, 2, 3, 4}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Generate(5) = %v, want %v", got, want)
	}
}

func TestGenerateEmpty(t *testing.T) {
	for _, n := range []int{0, -3} {
		got := Collect(Generate(n))
		if got == nil || len(got) != 0 {
			t.Fatalf("Generate(%d) = %v, want empty slice", n, got)
		}
	}
}

func Collect(in <-chan int) []int {
	out := make([]int, 0)
	for v := range in {
		out = append(out, v)
	}
	return out
}
