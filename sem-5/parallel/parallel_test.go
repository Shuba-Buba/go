package parallel

import (
	"errors"
	"sync/atomic"
	"testing"
)

func TestParallelSuccess(t *testing.T) {
	var n atomic.Int32
	fns := make([]func() error, 4)
	for i := range fns {
		fns[i] = func() error {
			n.Add(1)
			return nil
		}
	}
	if err := Parallel(fns...); err != nil || n.Load() != 4 {
		t.Fatal(err, n.Load())
	}
}

func TestParallelEmpty(t *testing.T) {
	if err := Parallel(); err != nil {
		t.Fatal(err)
	}
}

func TestParallelError(t *testing.T) {
	boom := errors.New("boom")
	err := Parallel(
		func() error { return nil },
		func() error { return boom },
		func() error { return nil },
	)
	if !errors.Is(err, boom) {
		t.Fatal(err)
	}
}
