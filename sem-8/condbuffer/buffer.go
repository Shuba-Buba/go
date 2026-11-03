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
	notEmpty *sync.Cond
	notFull  *sync.Cond

	data []any
	head int
	tail int
	size int
}

func New(capacity int) Buffer {
	b := &buffer{
		data: make([]any, capacity),
	}

	b.notEmpty = sync.NewCond(&b.mu)
	b.notFull = sync.NewCond(&b.mu)

	return b
}

func (b *buffer) Put(x any) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for b.size == len(b.data) {
		b.notFull.Wait()
	}

	b.data[b.tail] = x
	b.tail = (b.tail + 1) % len(b.data)
	b.size++

	b.notEmpty.Signal()
}

func (b *buffer) Get() any {
	b.mu.Lock()
	defer b.mu.Unlock()

	for b.size == 0 {
		b.notEmpty.Wait()
	}

	x := b.data[b.head]
	b.data[b.head] = nil

	b.head = (b.head + 1) % len(b.data)
	b.size--

	b.notFull.Signal()

	return x
}

func (b *buffer) Len() int {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.size
}
