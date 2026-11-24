//go:build !solution

package catalog

// IndexByID строит индекс при условии, что все идентификаторы уникальны.
func IndexByID[K comparable, T Identified[K]](records []T) (map[K]T, bool) {
	_ = records

	panic("TODO: implement IndexByID")
}
