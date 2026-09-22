//go:build !solution

package timeline

import "sort"

func Build(initial int, changes []Change) ([]Snapshot, error) {
	if initial < 0 {
		return nil, ErrNegativeValue
	}

	sums := make(map[int]int, len(changes))
	for _, c := range changes {
		sums[c.Day] += c.Delta
	}

	days := make([]int, 0, len(sums))
	for day := range sums {
		days = append(days, day)
	}
	sort.Ints(days)

	value := initial
	result := make([]Snapshot, 0, len(days))

	for _, day := range days {
		delta := sums[day]
		if delta == 0 {
			continue
		}

		value += delta
		if value < 0 {
			return nil, ErrNegativeValue
		}

		result = append(result, Snapshot{Day: day, Value: value})
	}

	return result, nil
}
