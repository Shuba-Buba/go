package linkstore

import (
	"encoding/json"
	"mime"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"sync"
	"testing"
)

const target = "https://example.org/a%2Fb?q=go#http"
const payload = `{"url":"` + target + `"}`

var validID = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

func request(h http.Handler, method, path, contentType, body string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, path, strings.NewReader(body))
	if contentType != "" {
		r.Header.Set("Content-Type", contentType)
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}

func createdLink(t *testing.T, w *httptest.ResponseRecorder, status int) string {
	t.Helper()

	if w.Code != status {
		t.Fatalf("POST status = %d, want %d; body = %q", w.Code, status, w.Body.String())
	}
	header := w.Result().Header
	mediaType, _, err := mime.ParseMediaType(header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		t.Fatalf("Content-Type = %q, want application/json", header.Get("Content-Type"))
	}
	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid response JSON: %v", err)
	}
	if len(body) != 1 || !validID.MatchString(body["id"]) {
		t.Fatalf("invalid response: %v, want an object with one valid id", body)
	}
	want := "/l/" + body["id"]
	if got := header.Get("Location"); got != want {
		t.Fatalf("Location = %q, want %q", got, want)
	}
	return want
}

func TestCreateAndRedirect(t *testing.T) {
	h := New()
	short := createdLink(t, request(h, "POST", "/links", "application/json; charset=utf-8", payload), 201)
	again := createdLink(t, request(h, "POST", "/links", "application/json", payload+"\n "), 200)
	if again != short {
		t.Fatalf("repeat created %q, want %q", again, short)
	}

	const otherTarget = "https://example.org/a%2fb?q=go#http"
	other := createdLink(t, request(h, "POST", "/links", "application/json", `{"url":"`+otherTarget+`"}`), 201)
	if other == short {
		t.Fatalf("two addresses got the same short link %q", short)
	}

	for _, method := range []string{"GET", "HEAD"} {
		w := request(h, method, short, "", "")
		if got := w.Result().Header.Get("Location"); w.Code != http.StatusFound || got != target {
			t.Errorf("%s %s = %d, Location %q; want 302, %q", method, short, w.Code, got, target)
		}
	}
	w := request(h, "GET", other, "", "")
	if got := w.Result().Header.Get("Location"); w.Code != http.StatusFound || got != otherTarget {
		t.Errorf("GET %s = %d, Location %q; want 302, %q", other, w.Code, got, otherTarget)
	}

	isolated := New()
	if w := request(isolated, "GET", short, "", ""); w.Code != http.StatusNotFound {
		t.Errorf("separate store: GET %s = %d, want 404", short, w.Code)
	}
}

func TestRouting(t *testing.T) {
	h := New()
	for _, tc := range []struct {
		method, path string
		status       int
		allow        []string
	}{
		{"GET", "/missing", 404, nil},
		{"GET", "/l/unknown", 404, nil},
		{"GET", "/l/two/segments", 404, nil},
		{"GET", "/links", 405, []string{"POST"}},
		{"DELETE", "/l/unknown", 405, []string{"GET", "HEAD"}},
	} {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			w := request(h, tc.method, tc.path, "", "")
			if w.Code != tc.status {
				t.Fatalf("status = %d, want %d", w.Code, tc.status)
			}
			if tc.allow == nil {
				return
			}
			allowed := map[string]bool{}
			for _, field := range w.Result().Header.Values("Allow") {
				for _, method := range strings.Split(field, ",") {
					allowed[strings.TrimSpace(method)] = true
				}
			}
			if len(allowed) != len(tc.allow) {
				t.Fatalf("Allow = %v, want %v", allowed, tc.allow)
			}
			for _, method := range tc.allow {
				if !allowed[method] {
					t.Errorf("Allow missing %s: %v", method, allowed)
				}
			}
		})
	}
}

func TestInvalidPost(t *testing.T) {
	h := New()
	for _, tc := range []struct {
		name, mediaType, body string
		status                int
	}{
		{"missing type", "", payload, 415},
		{"wrong type", "text/plain", payload, 415},
		{"invalid type", "application/json; x", payload, 415},
		{"empty", "application/json", "", 400},
		{"broken JSON", "application/json", `{"url":`, 400},
		{"extra field", "application/json", `{"url":"` + target + `","extra":1}`, 400},
		{"two values", "application/json", payload + " {}", 400},
		{"junk after value", "application/json", payload + " !", 400},
		{"missing URL", "application/json", "{}", 400},
		{"null", "application/json", "null", 400},
		{"empty URL", "application/json", `{"url":""}`, 400},
		{"relative", "application/json", `{"url":"/hello"}`, 400},
		{"wrong scheme", "application/json", `{"url":"ftp://example.org/x"}`, 400},
		{"no hostname", "application/json", `{"url":"https://:443/x"}`, 400},
		{"credentials", "application/json", `{"url":"https://user@example.org"}`, 400},
	} {
		t.Run(tc.name, func(t *testing.T) {
			w := request(h, "POST", "/links", tc.mediaType, tc.body)
			if w.Code != tc.status {
				t.Errorf("status = %d, want %d", w.Code, tc.status)
			}
		})
	}
	createdLink(t, request(h, "POST", "/links", "application/json", payload), 201)
}

func TestBodyLimitWithoutContentLength(t *testing.T) {
	h := New()
	for _, body := range []string{
		payload + strings.Repeat(" ", 1025-len(payload)),
		"!" + strings.Repeat(" ", 1024),
	} {
		r := httptest.NewRequest("POST", "/links", strings.NewReader(body))
		r.ContentLength = -1
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Code != http.StatusRequestEntityTooLarge {
			t.Errorf("status = %d for %d bytes, want 413", w.Code, len(body))
		}
	}
	exact := payload + strings.Repeat(" ", 1024-len(payload))
	createdLink(t, request(h, "POST", "/links", "application/json", exact), 201)
}

func TestConcurrentPosts(t *testing.T) {
	const n = 24
	h := New()
	start := make(chan struct{})
	out := make([]*httptest.ResponseRecorder, n)
	var wg sync.WaitGroup
	for i := range out {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			out[i] = request(h, "POST", "/links", "application/json", payload)
		}()
	}
	close(start)
	wg.Wait()

	created := 0
	for i, w := range out {
		switch w.Code {
		case http.StatusCreated:
			created++
		case http.StatusOK:
		default:
			t.Fatalf("POST %d: status = %d, want 200 or 201", i, w.Code)
		}
		link := createdLink(t, w, w.Code)
		if want := out[0].Result().Header.Get("Location"); link != want {
			t.Errorf("POST %d: Location = %q, want %q", i, link, want)
		}
	}
	if created != 1 {
		t.Errorf("%d POSTs got 201, want 1", created)
	}
}
