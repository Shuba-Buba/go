//go:build !solution

package catalog

// IndexByID строит индекс при условии, что все идентификаторы уникальны.
func IndexByID[K comparable, T Identified[K]](records []T) (map[K]T, bool) {
	index := make(map[K]T, len(records))
	for _, record := range records {
		id := record.ID()
		if _, ok := index[id]; ok {
			return nil, false
		}
		index[id] = record
	}
	return index, true
}
