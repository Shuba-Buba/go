package oncecache

import (
	"sync/atomic"
	"testing"
	"time"
)

type getResult struct {
	value any
	ok    bool
}

func asyncGet(c Cache, key string) <-chan getResult {
	result := make(chan getResult, 1)

	go func() {
		value, ok := c.Get(key)
		result <- getResult{value: value, ok: ok}
	}()

	return result
}

func receive(t *testing.T, result <-chan getResult) getResult {
	t.Helper()

	select {
	case got := <-result:
		return got
	case <-time.After(2 * time.Second):
		t.Fatal("Get did not finish")
		return getResult{}
	}
}

func TestCacheKeysNilAndInstances(t *testing.T) {
	var calls atomic.Int32

	c := New(func(key string) any {
		calls.Add(1)

		if key == "nil" {
			return nil
		}

		return "v:" + key
	})
	for range 2 {
		for _, key := range []string{"a", "b", "nil", ""} {
			got := receive(t, asyncGet(c, key))

			var want any = "v:" + key
			if key == "nil" {
				want = nil
			}

			if !got.ok || got.value != want {
				t.Fatalf("Get(%q)=%v,%v, want=%v,true", key, got.value, got.ok, want)
			}
		}
	}

	if got := calls.Load(); got != 4 {
		t.Fatalf("loader calls=%d, want=4", got)
	}

	other := New(func(_ string) any { return "other" })

	got := receive(t, asyncGet(other, "a"))
	if !got.ok || got.value != "other" {
		t.Fatalf("different cache returned %v,%v", got.value, got.ok)
	}
}

func TestCacheWaitsAndOtherKeysProgress(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})

	released := false
	defer func() {
		if !released {
			close(release)
		}
	}()

	var calls atomic.Int32

	c := New(func(key string) any {
		if key == "slow" {
			if calls.Add(1) == 1 {
				close(started)
			}

			<-release
		}

		return "v:" + key
	})
	first := asyncGet(c, "slow")

	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatal("loader did not start")
	}

	const peers = 16

	ready := make(chan struct{}, peers)

	results := make(chan getResult, peers)
	for range peers {
		go func() {
			ready <- struct{}{}

			value, ok := c.Get("slow")
			results <- getResult{value: value, ok: ok}
		}()
	}

	for range peers {
		select {
		case <-ready:
		case <-time.After(2 * time.Second):
			t.Fatal("caller did not start")
		}
	}
	// Independent progress is required while the first loader is still blocked.
	fast := asyncGet(c, "fast")
	select {
	case early := <-results:
		t.Fatalf("Get returned before its loader finished: %v", early)
	case got := <-fast:
		if !got.ok || got.value != "v:fast" {
			t.Fatalf("other key returned %v", got)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("a slow key blocked an independent key")
	}
	// ready proves startup, not entry into Get. This is an additional check;
	// the final values and exact load count also cover late-scheduled callers.
	select {
	case early := <-results:
		t.Fatalf("Get returned before its loader finished: %v", early)
	default:
	}

	close(release)

	released = true

	for _, result := range []<-chan getResult{first, results} {
		count := 1
		if result == results {
			count = peers
		}

		for range count {
			got := receive(t, result)
			if !got.ok || got.value != "v:slow" {
				t.Fatalf("slow key returned %v", got)
			}
		}
	}

	if got := calls.Load(); got != 1 {
		t.Fatalf("slow loader calls=%d, want=1", got)
	}
}
