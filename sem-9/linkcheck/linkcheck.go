//go:build !solution

package linkcheck

import (
	"context"
	"net/http"
)

type Result struct {
	Status int
	Err    error
}

func Check(ctx context.Context, client *http.Client, urls []string, limit int) []Result {
	// TODO: schedule the requests and collect results.
	return make([]Result, len(urls))
}
