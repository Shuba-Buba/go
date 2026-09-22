//go:build !solution

package slicestack

// Push добавляет value на вершину стека.
func Push(stack []int, value int) []int {
	_ = value

	return stack
}

// Peek возвращает верхнее значение, не меняя стек.
func Peek(stack []int) (value int, err error) {
	_ = stack

	return 0, nil
}

// Pop возвращает верхнее значение и оставшийся стек.
func Pop(stack []int) (value int, rest []int, err error) {
	return 0, stack, nil
}
