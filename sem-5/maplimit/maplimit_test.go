package maplimit

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestMapLimitOrder(t *testing.T) {
	items := []int{1, 2, 3, 4}
	got, err := MapLimit(items, 2, func(x int) (int, error) {
		return x * 10, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	want := []int{10, 20, 30, 40}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("MapLimit = %v, want %v", got, want)
		}
	}
}

func TestMapLimitEmpty(t *testing.T) {
	got, err := MapLimit(nil, 2, func(x int) (int, error) { return x, nil })
	if err != nil || got == nil || len(got) != 0 {
		t.Fatal(got, err)
	}
}

func TestMapLimitConcurrency(t *testing.T) {
	var running atomic.Int32
	var maxRunning atomic.Int32

	_, err := MapLimit([]int{1, 2, 3, 4, 5, 6, 7, 8}, 3, func(x int) (int, error) {
		cur := running.Add(1)
		for {
			old := maxRunning.Load()
			if cur <= old || maxRunning.CompareAndSwap(old, cur) {
				break
			}
		}
		time.Sleep(30 * time.Millisecond)
		running.Add(-1)
		return x, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if maxRunning.Load() > 3 {
		t.Fatalf("max concurrent = %d, want <= 3", maxRunning.Load())
	}
}

func TestMapLimitError(t *testing.T) {
	boom := errors.New("boom")
	got, err := MapLimit([]int{1, 2, 3}, 2, func(x int) (int, error) {
		if x == 2 {
			return 0, boom
		}
		return x, nil
	})
	if got != nil || !errors.Is(err, boom) {
		t.Fatal(got, err)
	}
}
