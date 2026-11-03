package mutexcounter

import (
	"sync"
	"testing"
	"time"
)

func TestCounterSimpleAndIndependent(t *testing.T) {
	first, second := New(), New()
	first.Inc()
	first.Add(4)
	first.Add(-2)
	first.Add(0)

	if got := first.Load(); got != 3 {
		t.Fatalf("Load()=%d, want=3", got)
	}

	if got := second.Load(); got != 0 {
		t.Fatalf("new counter=%d, want=0", got)
	}
}

func TestCounterConcurrentReadersAndWriters(t *testing.T) {
	c := New()

	const writers, readers, iterations = 8, 8, 2000

	start := make(chan struct{})
	done := make(chan struct{})

	var wg sync.WaitGroup
	wg.Add(writers + readers)

	for range writers {
		go func() {
			defer wg.Done()

			<-start

			for range iterations {
				c.Inc()
				c.Add(1)
			}
		}()
	}

	for range readers {
		go func() {
			defer wg.Done()

			<-start

			previous := 0

			for range iterations {
				current := c.Load()
				if current < previous || current > 2*writers*iterations {
					t.Errorf("invalid counter change: %d -> %d", previous, current)
					return
				}

				previous = current
			}
		}()
	}

	go func() {
		wg.Wait()
		close(done)
	}()

	close(start)

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("counter methods did not finish")
	}

	if got := c.Load(); got != 2*writers*iterations {
		t.Fatalf("Load()=%d, want=%d", got, 2*writers*iterations)
	}
}

func BenchmarkCounterAdd(b *testing.B) {
	c := New()

	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Add(1)
		}
	})
}
