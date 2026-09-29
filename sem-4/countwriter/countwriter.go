//go:build !solution

package countwriter

import "io"

func New(dst io.Writer) *Writer               { return &Writer{dst: dst} }
func (w *Writer) Write(p []byte) (int, error) { return 0, nil }
func (w *Writer) BytesWritten() int64         { return 0 }
