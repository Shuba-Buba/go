//go:build !solution

package accesslog

import (
	"context"
	"net/http"
)

type Entry struct {
	ID     string
	Method string
	Path   string
	Status int
	Bytes  int
}

func Wrap(next http.Handler, log func(Entry)) http.Handler {
	// TODO: attach the ID and record the response.
	return next
}

func RequestID(ctx context.Context) string {
	return ""
}
