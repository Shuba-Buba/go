package fuzzcase

import (
	"errors"
	"testing"
)

func TestParsePort(t *testing.T) {
	for _, tc := range []struct {
		name, input string
		want        uint16
		wantError   bool
	}{
		{name: "http", input: "80", want: 80},
		{name: "leading zeros", input: "00080", want: 80},
		{name: "minimum", input: "1", want: 1},
		{name: "maximum", input: "65535", want: 65535},
		{name: "empty", input: "", wantError: true},
		{name: "zero", input: "0", wantError: true},
		{name: "too large", input: "65536", wantError: true},
		{name: "spaces", input: " 80 ", wantError: true},
		{name: "text", input: "http", wantError: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParsePort(tc.input)
			if got != tc.want {
				t.Errorf("ParsePort(%q) = %d, want %d", tc.input, got, tc.want)
			}

			if tc.wantError {
				if !errors.Is(err, ErrPort) {
					t.Errorf("ParsePort(%q) error = %v, want ErrPort", tc.input, err)
				}
			} else if err != nil {
				t.Errorf("ParsePort(%q) error = %v, want nil", tc.input, err)
			}
		})
	}
}
