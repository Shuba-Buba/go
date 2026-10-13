package testcleanup

import (
	"errors"
	"os"
	"testing"
)

func TestFileLifecycle(t *testing.T) {
	for _, tc := range []struct {
		name, content string
	}{
		{name: "text", content: "Hello, testing!"},
		{name: "empty", content: ""},
		{name: "unicode", content: "Гофер учится\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var path string

			t.Run("create and read", func(t *testing.T) {
				path = setupTempFile(t, tc.content)

				data, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}

				if string(data) != tc.content {
					t.Errorf("file contains %q, want %q", data, tc.content)
				}
			})

			if path == "" {
				t.Fatal("setupTempFile did not return a path")
			}

			if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
				t.Errorf("file %q still exists after subtest cleanup: %v", path, err)
			}
		})
	}
}

func TestIndependentFiles(t *testing.T) {
	var paths [2]string

	t.Run("files", func(t *testing.T) {
		for i, content := range []string{"first", "second"} {
			paths[i] = setupTempFile(t, content)

			data, err := os.ReadFile(paths[i])
			if err != nil {
				t.Fatal(err)
			}

			if string(data) != content {
				t.Errorf("file %d contains %q, want %q", i, data, content)
			}
		}

		if paths[0] == paths[1] {
			t.Errorf("two calls reused the same file %q", paths[0])
		}
	})

	for _, path := range paths {
		if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
			t.Errorf("file %q was not removed: %v", path, err)
		}
	}
}
