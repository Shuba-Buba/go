package printqueue

import "errors"

var (
	ErrUnavailable  = errors.New("printer or paper unavailable")
	ErrInvalidPages = errors.New("pages must be positive")
)
