//go:build !solution

package accesslog

import (
	"context"
	"net/http"
	"strconv"
	"sync/atomic"
)

type Entry struct {
	ID     string
	Method string
	Path   string
	Status int
	Bytes  int
}

type contextKey struct{}

type responseWriter struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (w *responseWriter) WriteHeader(status int) {
	if w.status != 0 {
		return
	}

	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *responseWriter) Write(p []byte) (int, error) {
	if w.status == 0 {
		w.WriteHeader(http.StatusOK)
	}

	n, err := w.ResponseWriter.Write(p)
	w.bytes += n
	return n, err
}

func Wrap(next http.Handler, log func(Entry)) http.Handler {
	var counter atomic.Uint64

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = strconv.FormatUint(counter.Add(1), 10)
		}

		ctx := context.WithValue(r.Context(), contextKey{}, id)
		r = r.WithContext(ctx)

		w.Header().Set("X-Request-ID", id)

		rw := &responseWriter{
			ResponseWriter: w,
		}

		next.ServeHTTP(rw, r)

		status := rw.status
		if status == 0 {
			status = http.StatusOK
		}

		log(Entry{
			ID:     id,
			Method: r.Method,
			Path:   r.URL.Path,
			Status: status,
			Bytes:  rw.bytes,
		})
	})
}

func RequestID(ctx context.Context) string {
	id, _ := ctx.Value(contextKey{}).(string)
	return id
}
