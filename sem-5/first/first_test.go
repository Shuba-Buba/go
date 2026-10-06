package first

import (
	"testing"
	"time"
)

func sendAfter(delay time.Duration, vals ...int) <-chan int {
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

func TestFirstPicksFastest(t *testing.T) {
	slow := sendAfter(200*time.Millisecond, 1)
	fast := sendAfter(10*time.Millisecond, 42)
	got, ok := First(slow, fast)
	if !ok || got != 42 {
		t.Fatalf("First = (%d, %v), want (42, true)", got, ok)
	}
}

func TestFirstSingle(t *testing.T) {
	got, ok := First(sendAfter(0, 7))
	if !ok || got != 7 {
		t.Fatal(got, ok)
	}
}

func TestFirstAllClosed(t *testing.T) {
	a := make(chan int)
	b := make(chan int)
	close(a)
	close(b)
	if got, ok := First(a, b); ok || got != 0 {
		t.Fatalf("First(closed...) = (%d, %v), want (0, false)", got, ok)
	}
}

func TestFirstEmptyList(t *testing.T) {
	if got, ok := First(); ok || got != 0 {
		t.Fatal(got, ok)
	}
}

func TestFirstDoesNotCloseInputs(t *testing.T) {
	fast := make(chan int, 1)
	slow := make(chan int, 1)
	fast <- 1
	slow <- 100
	got, ok := First(fast, slow)
	if !ok {
		t.Fatal("expected value")
	}
	switch got {
	case 1:
		if v, open := <-slow; !open || v != 100 {
			t.Fatalf("slow channel corrupted: (%d, %v)", v, open)
		}
	case 100:
		if v, open := <-fast; !open || v != 1 {
			t.Fatalf("fast channel corrupted: (%d, %v)", v, open)
		}
	default:
		t.Fatalf("unexpected value %d", got)
	}
}
