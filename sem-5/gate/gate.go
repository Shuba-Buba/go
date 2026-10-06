//go:build !solution

package gate

func NewGate(capacity int) *Gate {
	return &Gate{slots: make(chan struct{}, capacity)}
}

func (g *Gate) Acquire() {}

func (g *Gate) Release() {}
