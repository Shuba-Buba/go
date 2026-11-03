//go:build !solution

package condbuffer

type Buffer interface {
	Put(x any)
	Get() any
	Len() int
}

func New(capacity int) Buffer {
	_ = capacity

	panic("implement me")
}
