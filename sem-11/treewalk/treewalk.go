//go:build !solution

package treewalk

// VisitLeaves передаёт значения листьев в visit слева направо.
func VisitLeaves[T any](root *Tree[T], visit func(T)) {
	if root == nil {
		return
	}
	if root.Left == nil && root.Right == nil {
		visit(root.Value)
		return
	}
	VisitLeaves(root.Left, visit)
	VisitLeaves(root.Right, visit)
}
