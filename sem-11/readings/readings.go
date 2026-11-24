//go:build !solution

package readings

import "cmp"

// RecordHighs возвращает последовательные рекорды в порядке наблюдения.
func RecordHighs[S ~[]E, E cmp.Ordered](values S) S {
	result := make(S, 0, len(values))

	var max E

	for i, value := range values {
		if i == 0 || value > max {
			result = append(result, value)
			max = value
		}
	}

	return result
}
