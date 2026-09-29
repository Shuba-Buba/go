package countwriter

import "io"

type Writer struct {
	dst   io.Writer
	count int64
}
