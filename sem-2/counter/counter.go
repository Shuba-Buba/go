//go:build !solution

package counter

// New возвращает функцию, выдающую арифметическую последовательность.
func New(start, step int) func() int {
	_ = step
	cur := start

	return func() int {
		cur += step
		return cur
	}
}
