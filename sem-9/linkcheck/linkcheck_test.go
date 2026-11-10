package linkcheck

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
	"testing/synctest"
	"time"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

type trackedBody struct {
	io.Reader
	close func()
}

func (b *trackedBody) Close() error {
	b.close()
	return nil
}

func TestStatusOrderConcurrencyAndClose(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		const base = "https://example.test"
		urls := []string{base + "/slow", base + "/down", base + "/error", base + "/slow", base + "/missing"}
		want := []int{204, 503, 0, 204, 404}
		offline := errors.New("offline")

		var mu sync.Mutex
		calls, closed := map[string]int{}, map[string]int{}
		active, peak := 0, 0
		client := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			path := r.URL.Path
			if r.Method != http.MethodHead {
				t.Errorf("method = %s for %s, want HEAD", r.Method, path)
			}

			mu.Lock()
			calls[path]++
			active++
			peak = max(peak, active)
			mu.Unlock()

			if path == "/error" {
				mu.Lock()
				active--
				mu.Unlock()
				return nil, offline
			}

			status := 404
			switch path {
			case "/slow":
				time.Sleep(100 * time.Millisecond)
				status = 204
			case "/down":
				status = 503
			}
			return &http.Response{
				StatusCode: status,
				Body: &trackedBody{Reader: strings.NewReader(""), close: func() {
					// Keep the slot occupied until Close completes.
					time.Sleep(time.Millisecond)
					mu.Lock()
					closed[path]++
					active--
					mu.Unlock()
				}},
			}, nil
		})}

		got := Check(context.Background(), client, urls, 2)
		if len(got) != len(want) {
			t.Fatalf("got %d results, want %d", len(got), len(want))
		}
		for i, res := range got {
			if res.Status != want[i] {
				t.Errorf("[%d] status = %d, want %d", i, res.Status, want[i])
			}
			if want[i] == 0 {
				if !errors.Is(res.Err, offline) {
					t.Errorf("[%d] error = %v, want offline", i, res.Err)
				}
			} else if res.Err != nil {
				t.Errorf("[%d] unexpected error: %v", i, res.Err)
			}
		}

		mu.Lock()
		defer mu.Unlock()
		if peak != 2 || active != 0 {
			t.Errorf("peak active = %d, still active = %d; want 2 and 0", peak, active)
		}
		for path, count := range map[string]int{"/slow": 2, "/down": 1, "/error": 1, "/missing": 1} {
			if calls[path] != count {
				t.Errorf("%s called %d times, want %d", path, calls[path], count)
			}
			wantClosed := count
			if path == "/error" {
				wantClosed = 0
			}
			if closed[path] != wantClosed {
				t.Errorf("%s closed %d bodies, want %d", path, closed[path], wantClosed)
			}
		}
	})
}

func TestUsesClientsRedirectPolicy(t *testing.T) {
	var paths []string
	client := &http.Client{
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			paths = append(paths, r.URL.Path)
			if r.URL.Path == "/final" {
				return &http.Response{StatusCode: 200, Body: http.NoBody}, nil
			}
			return &http.Response{
				StatusCode: 302,
				Header:     http.Header{"Location": {"/final"}},
				Body:       http.NoBody,
			}, nil
		}),
	}
	got := Check(context.Background(), client, []string{"https://example.test/start"}, 1)
	if len(got) != 1 || got[0].Status != 302 || got[0].Err != nil {
		t.Errorf("results = %+v, want a 302 response without error", got)
	}
	if len(paths) != 1 || paths[0] != "/start" {
		t.Errorf("requested %v, want [/start]", paths)
	}
}

func TestCancellationKeepsCompletedResults(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		const base = "https://example.test"
		urls := []string{base + "/a", base + "/b", base + "/c", base + "/d"}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		started := make(chan string, len(urls))
		finished := make(chan struct{})
		var mu sync.Mutex
		var paths []string
		client := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			mu.Lock()
			paths = append(paths, r.URL.Path)
			n := len(paths)
			mu.Unlock()
			started <- base + r.URL.Path

			if n == 1 {
				return &http.Response{StatusCode: 200, Body: http.NoBody}, nil
			}

			select {
			case <-r.Context().Done():
				if n == 2 {
					close(finished)
				}
				return nil, r.Context().Err()
			case <-time.After(time.Second):
				return nil, errors.New("request context was not canceled")
			}
		})}

		done := make(chan []Result, 1)
		go func() { done <- Check(ctx, client, urls, 1) }()

		var completedURL string
		for i := range 2 {
			select {
			case u := <-started:
				if i == 0 {
					completedURL = u
				}
			case <-time.After(time.Second):
				t.Fatal("second request never started")
			}
		}
		cancel()

		var got []Result
		select {
		case got = <-done:
		case <-time.After(2 * time.Second):
			t.Fatal("Check did not finish after cancellation")
		}
		select {
		case <-finished:
		default:
			t.Error("Check returned without waiting for the canceled request")
		}
		mu.Lock()
		if len(paths) != 2 {
			t.Errorf("requested %v, want only two requests", paths)
		}
		mu.Unlock()

		if len(got) != len(urls) {
			t.Fatalf("got %d results, want %d", len(got), len(urls))
		}
		for i, r := range got {
			if urls[i] == completedURL {
				if r.Status != 200 || r.Err != nil {
					t.Errorf("[%d] completed result = %+v, want 200 without error", i, r)
				}
			} else if r.Status != 0 || !errors.Is(r.Err, context.Canceled) {
				t.Errorf("[%d] canceled result = %+v", i, r)
			}
		}
	})
}

func TestEmptyAndAlreadyCanceled(t *testing.T) {
	client := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		t.Errorf("unexpected request for %s", r.URL)
		return nil, errors.New("unexpected request")
	})}

	if got := Check(context.Background(), client, nil, 10); len(got) != 0 {
		t.Errorf("empty input returned %v", got)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	got := Check(ctx, client, []string{"https://example.test/one", "https://example.test/two"}, 10)
	if len(got) != 2 {
		t.Fatalf("got %d results, want 2", len(got))
	}
	for i, res := range got {
		if res.Status != 0 || !errors.Is(res.Err, context.Canceled) {
			t.Errorf("[%d] got %+v, want context.Canceled", i, res)
		}
	}
}
