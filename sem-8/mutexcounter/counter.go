//go:build !solution

package mutexcounter

import "sync"

type Counter interface {
	Inc()
	Add(delta int)
	Load() int
}

type counter struct {
	mu    sync.Mutex
	value int
}

func (c *counter) Add(delta int) {
	c.mu.Lock()
	c.value += delta
	c.mu.Unlock()
}

func (c *counter) Inc() {
	c.Add(1)
}

func (c *counter) Load() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.value
}

func New() Counter {
	return &counter{}
}
