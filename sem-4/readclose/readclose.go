//go:build !solution

package readclose

import "io"

func Load(open func() (io.ReadCloser, error)) ([]byte, error) { return nil, nil }
