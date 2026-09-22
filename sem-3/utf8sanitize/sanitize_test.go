package utf8sanitize

import (
	"testing"
	"unicode/utf8"
)

func TestSanitizeValidInput(t *testing.T) {
	input := "Go, Привет, \uFFFD!"

	got, invalid := Sanitize(input)
	if got != input || invalid != 0 {
		t.Fatalf("Sanitize(%q) = (%q, %d), want (%q, 0)", input, got, invalid, input)
	}
}

func TestSanitizeInvalidBytes(t *testing.T) {
	input := string([]byte{'A', 0xff, 'B', 0xc0, 0x80, 'C'})
	want := "A\uFFFDB\uFFFD\uFFFDC"

	got, invalid := Sanitize(input)
	if got != want || invalid != 3 {
		t.Fatalf("Sanitize returned (%q, %d), want (%q, 3)", got, invalid, want)
	}

	if !utf8.ValidString(got) {
		t.Fatalf("Sanitize returned invalid UTF-8: %q", got)
	}
}

func TestSanitizeTruncatedEncoding(t *testing.T) {
	input := string([]byte{0xe2, 0x82})

	got, invalid := Sanitize(input)
	if got != "\uFFFD\uFFFD" || invalid != 2 {
		t.Fatalf("Sanitize returned (%q, %d), want (%q, 2)", got, invalid, "\uFFFD\uFFFD")
	}
}
