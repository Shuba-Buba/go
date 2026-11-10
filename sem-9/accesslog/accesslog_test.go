package accesslog

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestRequestIDAndContext(t *testing.T) {
	type key struct{}
	parent := context.WithValue(context.Background(), key{}, "parent")
	ctx, cancel := context.WithCancel(parent)
	cancel()
	req := httptest.NewRequest("POST", "/links?debug=1", nil).WithContext(ctx)
	req.Header.Set("X-Request-ID", "from-client")

	finished := false
	var entries []Entry
	h := Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := RequestID(r.Context()); got != "from-client" {
			t.Errorf("RequestID = %q, want from-client", got)
		}
		if got := r.Context().Value(key{}); got != "parent" || !errors.Is(r.Context().Err(), context.Canceled) {
			t.Errorf("parent value = %v, context error = %v", got, r.Context().Err())
		}
		if got := w.Header().Get("X-Request-ID"); got != "from-client" {
			t.Errorf("response ID before write = %q, want from-client", got)
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusAccepted)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "hi ")
		_, _ = io.WriteString(w, "мир")
		finished = true
	}), func(e Entry) {
		if !finished {
			t.Error("log called before handler returned")
		}
		entries = append(entries, e)
	})

	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	const body = "hi мир"
	if w.Code != 202 || w.Body.String() != body || w.Result().Header.Get("Content-Type") != "text/plain" {
		t.Errorf("response = %d %q %v", w.Code, w.Body.String(), w.Result().Header)
	}
	want := Entry{ID: "from-client", Method: "POST", Path: "/links", Status: 202, Bytes: len(body)}
	if len(entries) != 1 || entries[0] != want {
		t.Errorf("entries = %+v, want [%+v]", entries, want)
	}
	if got := RequestID(req.Context()); got != "" {
		t.Errorf("original context ID = %q, want empty", got)
	}
}

func TestImplicitStatus(t *testing.T) {
	for _, body := range []string{"", "body"} {
		t.Run(fmt.Sprintf("bytes=%d", len(body)), func(t *testing.T) {
			var entries []Entry
			h := Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if body != "" {
					_, _ = w.Write([]byte(body))
					w.WriteHeader(http.StatusTeapot)
				}
			}), func(e Entry) { entries = append(entries, e) })

			w := httptest.NewRecorder()
			h.ServeHTTP(w, httptest.NewRequest("GET", "/", nil))
			resp := w.Result()
			id := resp.Header.Get("X-Request-ID")
			if resp.StatusCode != 200 || w.Body.String() != body || id == "" {
				t.Errorf("response = %d %q, ID %q; want 200 %q and a generated ID", resp.StatusCode, w.Body.String(), id, body)
			}
			want := Entry{ID: id, Method: "GET", Path: "/", Status: 200, Bytes: len(body)}
			if len(entries) != 1 || entries[0] != want {
				t.Errorf("entries = %+v, want [%+v]", entries, want)
			}
		})
	}
}

type partialWriter struct {
	*httptest.ResponseRecorder
}

func (w partialWriter) Write(p []byte) (int, error) {
	n, _ := w.ResponseRecorder.Write(p[:min(len(p), 2)])
	return n, io.ErrShortWrite
}

func TestPartialWrite(t *testing.T) {
	var entries []Entry
	h := Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n, err := w.Write([]byte("hello"))
		if n != 2 || err != io.ErrShortWrite {
			t.Errorf("Write = %d, %v; want 2, %v", n, err, io.ErrShortWrite)
		}
	}), func(e Entry) { entries = append(entries, e) })
	w := partialWriter{httptest.NewRecorder()}
	h.ServeHTTP(w, httptest.NewRequest("GET", "/", nil))

	if w.Body.String() != "he" {
		t.Errorf("body = %q, want he", w.Body.String())
	}
	if len(entries) != 1 || entries[0].Bytes != 2 || entries[0].Status != 200 {
		t.Errorf("entries = %+v, want one entry with 2 bytes and status 200", entries)
	}
}

func TestConcurrentGeneratedIDs(t *testing.T) {
	const n = 20
	var mu sync.Mutex
	var entries []Entry
	h := Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, RequestID(r.Context()))
	}), func(e Entry) {
		mu.Lock()
		defer mu.Unlock()
		entries = append(entries, e)
	})

	start := make(chan struct{})
	var wg sync.WaitGroup
	wantIDs := make(map[string]string, n)
	for i := range n {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			path := fmt.Sprintf("/l/%d", i)
			w := httptest.NewRecorder()
			h.ServeHTTP(w, httptest.NewRequest("GET", path, nil))
			id := w.Result().Header.Get("X-Request-ID")
			if id == "" || w.Body.String() != id {
				t.Errorf("%s: response ID = %q, context ID = %q", path, id, w.Body.String())
			}
			mu.Lock()
			wantIDs[path] = id
			mu.Unlock()
		}()
	}
	close(start)
	wg.Wait()

	if len(entries) != n {
		t.Fatalf("%d entries, want %d", len(entries), n)
	}
	seen := make(map[string]bool, n)
	for _, e := range entries {
		id, ok := wantIDs[e.Path]
		if !ok || id == "" || id != e.ID || seen[e.ID] || e.Bytes != len(id) || e.Status != 200 {
			t.Errorf("unexpected or duplicate entry: %+v", e)
		}
		seen[e.ID] = true
	}
}
