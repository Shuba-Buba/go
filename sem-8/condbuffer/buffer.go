//go:build !solution

package condbuffer

import "sync"

type Buffer interface {
	Put(x any)
	Get() any
	Len() int
}

type buffer struct {
	mu       sync.Mutex
	cond     *sync.Cond
	queue    []any
	capacity int
}

func (b *buffer) Put(x any) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for len(b.queue) >= b.capacity {
		b.cond.Wait()
	}
	b.queue = append(b.queue, x)
	b.cond.Broadcast()
}

func (b *buffer) Get() any {
	b.mu.Lock()
	defer b.mu.Unlock()
	for len(b.queue) == 0 {
		b.cond.Wait()
	}
	x := b.queue[0]
	b.queue[0] = nil
	b.queue = b.queue[1:]

	b.cond.Broadcast()

	return x
}

func (b *buffer) Len() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return len(b.queue)
}

func New(capacity int) Buffer {
	b := &buffer{
		queue:    make([]any, 0, capacity),
		capacity: capacity,
	}
	b.cond = sync.NewCond(&b.mu)
	return b
}
