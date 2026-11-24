package treewalk

// Tree хранит значение и ссылки на дочерние узлы.
type Tree[T any] struct {
	Value       T
	Left, Right *Tree[T]
}
