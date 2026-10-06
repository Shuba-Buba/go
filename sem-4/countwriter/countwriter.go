//go:build !solution

package countwriter

import "io"

func New(dst io.Writer) *Writer {
	return &Writer{dst: dst}
}

func (w *Writer) Write(p []byte) (int, error) {
	n, err := w.dst.Write(p)
	w.count += int64(n)
	if err != nil {
		return n, err
	}
	if n < len(p) {
		return n, io.ErrShortWrite
	}
	return n, nil
}

func (w *Writer) BytesWritten() int64 {
	return w.count
}
