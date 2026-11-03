//go:build !solution

package barrier

import "sync"

type Barrier interface {
	Wait()
	Done()
}

type barrier struct {
	mu           sync.Mutex
	cond         *sync.Cond
	participants int
	waiting      int
	generation   int
}

func New(n int) Barrier {
	b := &barrier{
		participants: n,
	}

	b.cond = sync.NewCond(&b.mu)
	return b
}

func (b *barrier) nextRound() {
	b.waiting = 0
	b.generation++
	b.cond.Broadcast()
}

func (b *barrier) Wait() {
	b.mu.Lock()
	defer b.mu.Unlock()

	gen := b.generation
	b.waiting++

	if b.waiting == b.participants {
		b.nextRound()
		return
	}

	for gen == b.generation {
		b.cond.Wait()
	}
}

func (b *barrier) Done() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.participants--

	if b.waiting == b.participants {
		b.nextRound()
	}
}
