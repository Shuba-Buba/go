//go:build !solution

package slicestack

// Push добавляет value на вершину стека.
func Push(stack []int, value int) []int {
	stack = append(stack, value)
	return stack
}

// Peek возвращает верхнее значение, не меняя стек.
func Peek(stack []int) (value int, err error) {
	if len(stack) == 0 {
		return 0, ErrEmpty
	}
	return stack[len(stack)-1], nil
}

// Pop возвращает верхнее значение и оставшийся стек.
func Pop(stack []int) (value int, rest []int, err error) {
	if len(stack) == 0 {
		return 0, stack, ErrEmpty
	}
	return stack[len(stack)-1], stack[:len(stack)-1], nil
}
