//go:build !solution

package testhelper

import (
	"strings"
)

func assertEqual(t testingT, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("assertEqual failed: %d and %d are not equal", got, want)
	}
}

func assertContains(t testingT, s, substr string) {
	t.Helper()
	if !strings.Contains(s, substr) {
		t.Errorf("assertContains faild: %q is not in %q", s, substr)
	}
}
