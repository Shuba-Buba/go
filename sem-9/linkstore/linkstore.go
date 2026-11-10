//go:build !solution

package linkstore

import (
	"bytes"
	"encoding/json"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

type store struct {
	mu    sync.Mutex
	byURL map[string]string
	byID  map[string]string
	next  uint64
}

func New() http.Handler {
	s := &store{
		byURL: make(map[string]string),
		byID:  make(map[string]string),
	}

	return http.HandlerFunc(s.handle)
}

func (s *store) handle(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path == "/links":
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", "POST")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		s.create(w, r)

	case strings.HasPrefix(r.URL.Path, "/l/") &&
		len(r.URL.Path) > len("/l/") &&
		!strings.Contains(r.URL.Path[len("/l/"):], "/"):

		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.Header().Set("Allow", "GET, HEAD")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		s.redirect(w, r)

	default:
		http.NotFound(w, r)
	}
}

func (s *store) create(w http.ResponseWriter, r *http.Request) {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		http.Error(w, "unsupported media type", http.StatusUnsupportedMediaType)
		return
	}

	data, err := io.ReadAll(io.LimitReader(r.Body, 1025))
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if len(data) > 1024 {
		http.Error(w, "body too large", http.StatusRequestEntityTooLarge)
		return
	}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()

	var req struct {
		URL string `json:"url"`
	}

	if err := decoder.Decode(&req); err != nil || req.URL == "" {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	parsed, err := url.Parse(req.URL)
	if err != nil ||
		(parsed.Scheme != "http" && parsed.Scheme != "https") ||
		parsed.Hostname() == "" ||
		parsed.User != nil {
		http.Error(w, "invalid URL", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	id, exists := s.byURL[req.URL]
	status := http.StatusOK

	if !exists {
		s.next++
		id = strconv.FormatUint(s.next, 10)
		s.byURL[req.URL] = id
		s.byID[id] = req.URL
		status = http.StatusCreated
	}
	s.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", "/l/"+id)
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(struct {
		ID string `json:"id"`
	}{ID: id})
}

func (s *store) redirect(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/l/")

	s.mu.Lock()
	target, ok := s.byID[id]
	s.mu.Unlock()

	if !ok {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Location", target)
	w.WriteHeader(http.StatusFound)
}
