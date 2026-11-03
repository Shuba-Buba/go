//go:build !solution

package mutexcounter

type Counter interface {
	Inc()
	Add(delta int)
	Load() int
}

func New() Counter {
	panic("implement me")
}
