//go:build !solution

package gauge

import (
	"fmt"
)

func (g *Gauge) Add(delta int) error {
	if g.level+delta < 0 {
		return ErrNegative
	}
	g.level += delta
	return nil
}

func (g Gauge) Value() int {
	return g.level
}

func (g Gauge) String() string {
	return fmt.Sprintf("level=%d", g.level)
}
