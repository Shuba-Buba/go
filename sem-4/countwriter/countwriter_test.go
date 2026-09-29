package countwriter

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

type short struct {
	calls int
	err   error
}

func (s *short) Write(p []byte) (int, error) { s.calls++; return 2, s.err }

func TestBufferAndInput(t *testing.T) {
	var dst bytes.Buffer
	w := New(&dst)
	var _ io.Writer = w
	p := []byte("Go")
	if n, err := w.Write(p); n != 2 || err != nil {
		t.Fatal(n, err)
	}
	if n, err := w.Write([]byte("!")); n != 1 || err != nil {
		t.Fatal(n, err)
	}
	if w.BytesWritten() != 3 || dst.String() != "Go!" || string(p) != "Go" {
		t.Fatal(w.BytesWritten(), dst.String(), p)
	}
}
func TestCountPartialFailure(t *testing.T) {
	broken := errors.New("disk")
	for _, err := range []error{broken, nil} {
		dst := &short{err: err}
		w := New(dst)
		if dst.calls != 0 {
			t.Fatal("eager write")
		}
		n, got := w.Write([]byte("hello"))
		want := err
		if want == nil {
			want = io.ErrShortWrite
		}
		if n != 2 || got != want || w.BytesWritten() != 2 || dst.calls != 1 {
			t.Fatal(n, got, w.BytesWritten(), dst.calls)
		}
	}
}
