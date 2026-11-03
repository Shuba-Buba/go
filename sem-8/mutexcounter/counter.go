
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

func New() Counter {
	return &counter{}
}

func (c *counter) Inc() {
	c.Add(1)
}

func (c *counter) Add(delta int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.value += delta
}

func (c *counter) Load() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.value
}
