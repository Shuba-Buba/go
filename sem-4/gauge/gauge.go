//go:build !solution

package gauge

func (g *Gauge) Add(delta int) error { return nil }
func (g Gauge) Value() int           { return 0 }
func (g Gauge) String() string       { return "" }
