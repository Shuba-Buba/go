//go:build !solution

package treewalk

// VisitLeaves передаёт значения листьев в visit слева направо.
func VisitLeaves[T any](root *Tree[T], visit func(T)) {
	_, _ = root, visit

	panic("TODO: implement VisitLeaves")
}
