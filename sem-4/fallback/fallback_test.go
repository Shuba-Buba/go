package fallback

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

type source struct {
	value string
	err   error
	keys  []string
}

func (s *source) Lookup(key string) (string, error) {
	s.keys = append(s.keys, key)
	return s.value, s.err
}

type denied struct{}

func (*denied) Error() string { return "denied" }

func TestOnlyMissingUsesBackup(t *testing.T) {
	for _, value := range []string{"cached", ""} {
		p, b := &source{value: value}, &source{value: "remote"}
		got, err := Lookup(p, b, "key")
		if err != nil || got != value || len(b.keys) != 0 || len(p.keys) != 1 {
			t.Fatal(got, err)
		}
	}
	p, b := &source{err: fmt.Errorf("cache: %w", ErrMissing)}, &source{value: "remote"}
	if got, err := Lookup(p, b, "key"); got != "remote" || err != nil || len(b.keys) != 1 || b.keys[0] != "key" {
		t.Fatal(got, err)
	}
}
func TestErrorChain(t *testing.T) {
	cause := &denied{}
	p, b := &source{err: cause}, &source{}
	_, err := Lookup(p, b, "key")
	var detail *denied
	if !errors.As(err, &detail) || detail != cause || !strings.Contains(err.Error(), "primary") || len(b.keys) != 0 {
		t.Fatal(err)
	}
	p.err = ErrMissing
	b.err = fmt.Errorf("database: %w", cause)
	_, err = Lookup(p, b, "key")
	if !errors.Is(err, cause) || !strings.Contains(err.Error(), "backup") {
		t.Fatal(err)
	}
}
