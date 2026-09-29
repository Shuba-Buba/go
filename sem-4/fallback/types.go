package fallback

import "errors"

var ErrMissing = errors.New("missing")

type Source interface {
	Lookup(key string) (string, error)
}
