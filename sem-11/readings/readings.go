//go:build !solution

package readings

import "cmp"

// RecordHighs возвращает последовательные рекорды в порядке наблюдения.
func RecordHighs[S ~[]E, E cmp.Ordered](values S) S {
	if len(values) == 0 {
		return nil
	}
	res := make(S, 1, len(values))
	res[0] = values[0]
	m_val := values[0]
	for _, v := range values[1:] {
		if v > m_val {
			m_val = v
			res = append(res, v)
		}
	}
	return res
}
