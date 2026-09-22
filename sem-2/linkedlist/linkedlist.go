//go:build !solution

package linkedlist

// Prepend добавляет новый узел в начало списка.
func Prepend(head *Node, value int) *Node {
	return &Node{Value: value, Next: head}
}

// Length возвращает число узлов списка.
func Length(head *Node) int {
	count := 0
	for cur := head; cur != nil; cur = cur.Next {
		count++
	}
	return count
}

// Find ищет первый узел со значением value.
func Find(head *Node, value int) *Node {
	for cur := head; cur != nil; cur = cur.Next {
		if cur.Value == value {
			return cur
		}
	}
	return nil
}

// Reverse разворачивает список, меняя связи между существующими узлами.
func Reverse(head *Node) *Node {
	var prev *Node
	cur := head
	for cur != nil {
		next := cur.Next
		cur.Next = prev
		prev = cur
		cur = next
	}

	return prev
}
