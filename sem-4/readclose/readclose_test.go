package readclose

import (
	"errors"
	"io"
	"strings"
	"testing"
)

type resource struct {
	reader            io.Reader
	readErr, closeErr error
	reads, closes     int
}

func (r *resource) Read(p []byte) (int, error) {
	r.reads++
	if r.readErr != nil {
		return 0, r.readErr
	}
	return r.reader.Read(p)
}
func (r *resource) Close() error { r.closes++; return r.closeErr }

func TestLifecycle(t *testing.T) {
	brokenRead, brokenClose := errors.New("read failure"), errors.New("close failure")
	for _, tc := range []struct {
		readErr, closeErr, want error
		stage                   string
	}{
		{nil, nil, nil, ""}, {brokenRead, nil, brokenRead, "read"}, {nil, brokenClose, brokenClose, "close"}, {brokenRead, brokenClose, brokenRead, "read"},
	} {
		r := &resource{reader: strings.NewReader("Go"), readErr: tc.readErr, closeErr: tc.closeErr}
		opens := 0
		data, err := Load(func() (io.ReadCloser, error) { opens++; return r, nil })
		if r.closes != 1 || opens != 1 {
			t.Fatal(opens, r.closes)
		}
		if tc.want == nil {
			if err != nil || string(data) != "Go" {
				t.Fatal(string(data), err)
			}
		} else {
			if data != nil || !errors.Is(err, tc.want) || !strings.Contains(err.Error(), tc.stage) {
				t.Fatal(data, err)
			}
		}
	}
}
func TestOpenFailure(t *testing.T) {
	cause := errors.New("no file")
	r := &resource{}
	data, err := Load(func() (io.ReadCloser, error) { return r, cause })
	if data != nil || !errors.Is(err, cause) || !strings.Contains(err.Error(), "open") || r.reads != 0 || r.closes != 0 {
		t.Fatal(data, err, r)
	}
}
