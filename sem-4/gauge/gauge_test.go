package gauge

import (
	"errors"
	"fmt"
	"testing"
)

func TestStateAndCopy(t *testing.T) {
	var g Gauge
	if g.Value() != 0 || g.String() != "level=0" {
		t.Fatal(g)
	}
	if err := g.Add(7); err != nil {
		t.Fatal(err)
	}
	copy := g
	if err := g.Add(-2); err != nil {
		t.Fatal(err)
	}
	if g.Value() != 5 || copy.Value() != 7 {
		t.Fatal(g, copy)
	}
	if err := g.Add(-6); !errors.Is(err, ErrNegative) || g.Value() != 5 {
		t.Fatal(err, g)
	}
	var view fmt.Stringer = g
	if view.String() != "level=5" {
		t.Fatal(view)
	}
}
