//go:build !solution

package utf8sanitize

import (
	"strings"
	"unicode/utf8"
)

func Sanitize(input string) (normalized string, invalidBytes int) {
	var b strings.Builder
	b.Grow(len(input))

	for i := 0; i < len(input); {
		r, size := utf8.DecodeRuneInString(input[i:])
		if r == utf8.RuneError && size == 1 {
			b.WriteRune(utf8.RuneError)
			invalidBytes++
			i++
			continue
		}

		b.WriteRune(r)
		i += size
	}

	return b.String(), invalidBytes
}
