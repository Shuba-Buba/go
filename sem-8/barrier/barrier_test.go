package barrier

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestBarrier_SyncTwoRounds(t *testing.T) {
	b := New(3)

	const rounds = 5

	var wg sync.WaitGroup
	wg.Add(3)

	for i := 0; i < 3; i++ {
		go func() {
			defer wg.Done()

			for r := 0; r < rounds; r++ {
				time.Sleep(time.Millisecond) // работа
				b.Wait()
			}
		}()
	}

	done := make(chan struct{})

	go func() { wg.Wait(); close(done) }()

	select {
	case <-done:
		// ok
	case <-time.After(3 * time.Second):
		t.Fatal("barrier deadlock")
	}
}

// Проверяем строгую цикличность: на каждом раунде должны собираться ровно N участников,
// счётчик не должен «перетекать» между раундами.
func TestBarrier_CyclicStrict(t *testing.T) {
	const (
		N      = 4
		rounds = 50
	)

	b := New(N)

	var arrived [rounds]atomic.Int32

	var wg sync.WaitGroup
	wg.Add(N)

	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()

			for r := 0; r < rounds; r++ {
				// небольшой джиттер прихода
				time.Sleep(time.Duration(rand.Intn(3)) * time.Millisecond)

				n := arrived[r].Add(1)
				if n > N {
					t.Errorf("arrived overflow: round=%d have=%d", r, n)
					return
				}

				b.Wait()

				// после выхода со встречи должны быть все N
				if arrived[r].Load() != N {
					t.Errorf(
						"barrier didn't wait: round=%d have=%d want=%d",
						r,
						arrived[r].Load(),
						N,
					)

					return
				}
			}
		}()
	}

	done := make(chan struct{})

	go func() { wg.Wait(); close(done) }()

	select {
	case <-done:
		// ok
	case <-time.After(5 * time.Second):
		t.Fatal("barrier deadlock or too slow")
	}
}

// Проверяем, что Done уменьшает кворум для следующих раундов,
// и текущий раунд правильно выпускается.
func TestBarrier_DoneCyclic(t *testing.T) {
	b := New(3)

	var (
		round0 atomic.Int32
		round1 atomic.Int32
	)

	var wg sync.WaitGroup
	wg.Add(3)

	// Путник 1 уходит после первой точки
	go func() {
		defer wg.Done()

		time.Sleep(5 * time.Millisecond)
		round0.Add(1)
		b.Wait()
		b.Done()
	}()

	// Путник 2
	go func() {
		defer wg.Done()

		round0.Add(1)
		b.Wait()

		time.Sleep(2 * time.Millisecond)
		round1.Add(1)
		b.Wait()

		if round1.Load() != 2 { // на второй точке должны встретиться вдвоём
			t.Errorf("second round not synchronized: have=%d want=2", round1.Load())
		}
	}()

	// Путник 3
	go func() {
		defer wg.Done()

		round0.Add(1)
		b.Wait()

		round1.Add(1)
		b.Wait()
	}()

	done := make(chan struct{})

	go func() { wg.Wait(); close(done) }()

	select {
	case <-done:
		// ok
	case <-time.After(3 * time.Second):
		t.Fatal("barrier Done() stuck")
	}

	if round0.Load() != 3 {
		t.Fatalf("first round arrivals=%d, want=3", round0.Load())
	}

	if round1.Load() != 2 {
		t.Fatalf("second round arrivals=%d, want=2", round1.Load())
	}
}

func TestBarrier_DoneShrinks(t *testing.T) {
	b := New(3)

	var wg sync.WaitGroup
	wg.Add(3)

	go func() { defer wg.Done(); b.Wait(); b.Done() }()
	go func() { defer wg.Done(); b.Wait(); b.Wait() }()
	go func() { defer wg.Done(); b.Wait(); b.Wait() }()

	done := make(chan struct{})

	go func() { wg.Wait(); close(done) }()

	select {
	case <-done:
		// ok
	case <-time.After(3 * time.Second):
		t.Fatal("barrier Done() stuck")
	}
}

func TestBarrierDoneBeforeFirstRound(t *testing.T) {
	b := New(3)
	b.Done()

	var arrived atomic.Int32

	done := make(chan struct{}, 2)

	for range 2 {
		go func() {
			arrived.Add(1)
			b.Wait()

			if got := arrived.Load(); got != 2 {
				t.Errorf("barrier released before both participants arrived: %d", got)
			}

			done <- struct{}{}
		}()
	}

	timer := time.NewTimer(2 * time.Second)
	defer timer.Stop()

	for range 2 {
		select {
		case <-done:
		case <-timer.C:
			t.Fatal("Done before the first round did not shrink the barrier")
		}
	}
}
