//go:build !solution

package fuzzcase

import "strconv"

// ParsePort parses a decimal port in the range 1 through 65535.
func ParsePort(s string) (uint16, error) {
	n, err := strconv.Atoi(s)
	if err != nil || n < 1 || n > 65535 {
		return 0, ErrPort
	}

	return uint16(n), nil
}
