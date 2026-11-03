//go:build !solution

package barrier

type Barrier interface {
	Wait()
	Done()
}

func New(n int) Barrier {
	_ = n

	panic("implement me")
}
