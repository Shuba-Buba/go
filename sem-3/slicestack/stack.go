//go:build !solution

package slicestack

func Push(stack []int, value int) []int {
	return append(stack, value)
}

func Peek(stack []int) (value int, err error) {
	if len(stack) == 0 {
		return 0, ErrEmpty
	}

	return stack[len(stack)-1], nil
}

func Pop(stack []int) (value int, rest []int, err error) {
	if len(stack) == 0 {
		return 0, stack, ErrEmpty
	}

	n := len(stack)

	return stack[n-1], stack[:n-1], nil
}
