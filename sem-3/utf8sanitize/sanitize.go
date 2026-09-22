//go:build !solution

package utf8sanitize

// Sanitize заменяет невалидные bytes на U+FFFD.
func Sanitize(input string) (normalized string, invalidBytes int) {
	return input, 0
}
