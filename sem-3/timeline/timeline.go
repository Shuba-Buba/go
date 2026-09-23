//go:build !solution

package timeline

import "slices"

// Build группирует изменения и строит накопительный timeline.
func Build(initial int, changes []Change) ([]Snapshot, error) {
	if initial < 0 {
		return nil, ErrNegativeValue
	}
	if len(changes) == 0 {
		return make([]Snapshot, 0), nil
	}

	changes_sorted := make([]Change, len(changes))
	for i, c := range changes {
		changes_sorted[i] = c
	}
	slices.SortFunc(changes_sorted, func(a, b Change) int {
		return a.Day - b.Day
	})
	result := make([]Snapshot, 0, len(changes_sorted))
	current := initial
	current_day := changes_sorted[0].Day
	for _, change := range changes_sorted {
		if change.Day != current_day {
			if current != initial {
				if current < 0 {
					return nil, ErrNegativeValue
				}
				result = append(result, Snapshot{current_day, current})
			}
			initial = current
			current_day = change.Day
		}
		current += change.Delta
		//if current < 0 {
		//	return nil, ErrNegativeValue
		//}
	}
	if current != initial {
		if current < 0 {
			return nil, ErrNegativeValue
		}
		result = append(result, Snapshot{current_day, current})
	}
	return result, nil
}
