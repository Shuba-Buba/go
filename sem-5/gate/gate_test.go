package gate

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestGateLimitsConcurrency(t *testing.T) {
	g := NewGate(2)
	var running atomic.Int32
	var maxRunning atomic.Int32
	var wg sync.WaitGroup

	for i := 0; i < 6; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			g.Acquire()
			defer g.Release()
			cur := running.Add(1)
			for {
				old := maxRunning.Load()
				if cur <= old || maxRunning.CompareAndSwap(old, cur) {
					break
				}
			}
			time.Sleep(20 * time.Millisecond)
			running.Add(-1)
		}()
	}
	wg.Wait()
	if maxRunning.Load() > 2 {
		t.Fatalf("max concurrent = %d, want <= 2", maxRunning.Load())
	}
}
