//go:build !solution

package timeline

// Build группирует изменения и строит накопительный timeline.
func Build(initial int, changes []Change) ([]Snapshot, error) {
	_, _ = initial, changes

	return nil, nil
}
