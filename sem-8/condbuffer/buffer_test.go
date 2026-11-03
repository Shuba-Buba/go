package condbuffer

import (
	"sync"
	"testing"
	"time"
)

func await(t *testing.T, done <-chan struct{}) {
	t.Helper()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("buffer operation did not finish")
	}
}

func get(t *testing.T, b Buffer) any {
	t.Helper()

	result := make(chan any, 1)
	go func() { result <- b.Get() }()

	select {
	case value := <-result:
		return value
	case <-time.After(2 * time.Second):
		t.Fatal("Get did not finish")
		return nil
	}
}

func put(t *testing.T, b Buffer, value any) {
	t.Helper()

	done := make(chan struct{})

	go func() {
		b.Put(value)
		close(done)
	}()

	await(t, done)
}

func TestBufferFIFOAndValues(t *testing.T) {
	b := New(4)

	values := []any{nil, 0, 7, 7}
	for i, value := range values {
		put(t, b, value)

		if got := b.Len(); got != i+1 {
			t.Fatalf("Len()=%d, want=%d", got, i+1)
		}
	}

	for i, want := range values {
		if got := get(t, b); got != want {
			t.Fatalf("Get()=%v, want=%v", got, want)
		}

		if got := b.Len(); got != len(values)-i-1 {
			t.Fatalf("Len()=%d after Get", got)
		}
	}

	other := New(1)
	put(t, other, "other")

	if got := b.Len(); got != 0 {
		t.Fatalf("different buffer changed the first: Len=%d", got)
	}

	if got := get(t, other); got != "other" {
		t.Fatalf("different buffer returned %v", got)
	}
}

func TestBufferFullAndEmpty(t *testing.T) {
	b := New(1)
	put(t, b, 1)

	started, done := make(chan struct{}), make(chan struct{})
	go func() {
		close(started)
		b.Put(2)
		close(done)
	}()

	await(t, started)
	// Startup is not proof that Put has entered: this is an additional check.
	select {
	case <-done:
		t.Fatal("Put returned while buffer was full")
	default:
	}

	if got := get(t, b); got != 1 {
		t.Fatalf("Get()=%v, want=1", got)
	}

	await(t, done)

	if got := get(t, b); got != 2 {
		t.Fatalf("Get()=%v, want=2", got)
	}

	result := make(chan any, 1)

	started = make(chan struct{})
	go func() {
		close(started)

		result <- b.Get()
	}()

	await(t, started)

	select {
	case got := <-result:
		t.Fatalf("Get returned from an empty buffer: %v", got)
	default:
	}

	put(t, b, 3)

	select {
	case got := <-result:
		if got != 3 {
			t.Fatalf("Get()=%v, want=3", got)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Get did not wake after Put")
	}
}

func TestBufferParallelKeepsEveryValue(t *testing.T) {
	const (
		workers, perWorker = 8, 200
		total              = workers * perWorker
	)

	b := New(8)
	values := make(chan any, total)
	start, done := make(chan struct{}), make(chan struct{})

	var wg sync.WaitGroup
	wg.Add(2 * workers)

	for producer := range workers {
		go func() {
			defer wg.Done()

			<-start

			for i := range perWorker {
				b.Put(producer*perWorker + i)
			}
		}()
		go func() {
			defer wg.Done()

			<-start

			for range perWorker {
				values <- b.Get()

				if length := b.Len(); length < 0 || length > 8 {
					t.Errorf("invalid concurrent Len=%d", length)
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(values)
		close(done)
	}()

	close(start)
	await(t, done)

	seen := make([]bool, total)
	count := 0

	for value := range values {
		index, ok := value.(int)
		if !ok || index < 0 || index >= total {
			t.Fatalf("unexpected value %v", value)
		}

		if seen[index] {
			t.Fatalf("duplicate value %d", index)
		}

		seen[index] = true
		count++
	}

	if count != total || b.Len() != 0 {
		t.Fatalf("received=%d, remaining=%d", count, b.Len())
	}
}
