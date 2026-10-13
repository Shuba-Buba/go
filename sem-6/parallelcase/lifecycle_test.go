//go:build !solution

package parallelcase

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestParallelFixture(t *testing.T) {
	path := filepath.Join(t.TempDir(), "input.txt")
	events := make(chan string, 4)

	t.Run("group", func(t *testing.T) {
		if err := os.WriteFile(path, []byte("shared input"), 0o600); err != nil {
			t.Fatal(err)
		}

		// TODO: сохраните ресурс до завершения обоих параллельных подтестов.
		t.Cleanup(func() {
			if err := os.Remove(path); err != nil {
				t.Error(err)
			}

			events <- "cleanup"
		})

		for _, name := range []string{"first", "second"} {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				data, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}

				if string(data) != "shared input" {
					t.Errorf("input = %q, want shared input", data)
				}

				events <- name
			})
		}

		events <- "parent body"
	})
	close(events)

	var trace []string

	for event := range events {
		trace = append(trace, event)
	}

	if len(trace) != 4 {
		t.Fatalf("events = %q, want parent body, two readers, cleanup", trace)
	}

	if trace[0] != "parent body" || trace[3] != "cleanup" || trace[1] == trace[2] {
		t.Errorf("invalid lifecycle: %q", trace)
	}

	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("fixture was not removed after group completion: %v", err)
	}
}
