//go:build !solution

package oncecache

type Cache interface {
	Get(key string) (any, bool)
}

func New(load func(string) any) Cache {
	_ = load

	panic("implement me")
}
