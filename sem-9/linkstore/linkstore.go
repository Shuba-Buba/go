//go:build !solution

package linkstore

import "net/http"

func New() http.Handler {
	return http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		// TODO: register routes and store links.
	})
}
