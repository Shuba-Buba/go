package testhelper

import (
	"fmt"
	"strings"
	"testing"
)

type fakeT struct {
	helpers            int
	errors             []string
	helpersBeforeError []int
}

func (f *fakeT) Helper() {
	f.helpers++
}

func (f *fakeT) Errorf(format string, args ...any) {
	f.errors = append(f.errors, fmt.Sprintf(format, args...))
	f.helpersBeforeError = append(f.helpersBeforeError, f.helpers)
}

func checkReport(t *testing.T, f *fakeT, wantError bool, details ...string) {
	t.Helper()

	if f.helpers != 1 {
		t.Errorf("Helper called %d times, want 1", f.helpers)
	}

	if !wantError {
		if len(f.errors) != 0 {
			t.Errorf("matching values reported errors: %q", f.errors)
		}

		return
	}

	if len(f.errors) != 1 {
		t.Fatalf("reported %d errors, want 1: %q", len(f.errors), f.errors)
	}

	if f.helpersBeforeError[0] != 1 {
		t.Errorf("Helper must run before Errorf, got %d earlier calls", f.helpersBeforeError[0])
	}

	for _, detail := range details {
		if !strings.Contains(f.errors[0], detail) {
			t.Errorf("error %q does not describe %q", f.errors[0], detail)
		}
	}
}

func TestAssertEqual(t *testing.T) {
	for _, tc := range []struct {
		name      string
		got, want int
	}{
		{name: "positive", got: 42, want: 42},
		{name: "negative", got: -17, want: -17},
		{name: "zero", got: 0, want: 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var f fakeT

			assertEqual(&f, tc.got, tc.want)
			checkReport(t, &f, false)
		})
	}
}

func TestAssertEqualFailure(t *testing.T) {
	var f fakeT

	assertEqual(&f, 17, 42)
	checkReport(t, &f, true, "17", "42")
}

func TestAssertContains(t *testing.T) {
	for _, tc := range []struct {
		name, s, substr string
	}{
		{name: "word", s: "learn Go", substr: "Go"},
		{name: "whole string", s: "кот", substr: "кот"},
		{name: "unicode", s: "учимся в ВШЭ", substr: "ВШЭ"},
		{name: "empty substring", s: "Go", substr: ""},
		{name: "both empty", s: "", substr: ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var f fakeT

			assertContains(&f, tc.s, tc.substr)
			checkReport(t, &f, false)
		})
	}
}

func TestAssertContainsFailure(t *testing.T) {
	var f fakeT

	assertContains(&f, "learn Go", "Rust")
	checkReport(t, &f, true, "learn Go", "Rust")
}

func TestAssertContainsEmptyInputFailure(t *testing.T) {
	var f fakeT

	assertContains(&f, "", "Go")
	checkReport(t, &f, true, "Go")
}
