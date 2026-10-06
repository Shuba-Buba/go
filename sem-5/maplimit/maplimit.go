//go:build !solution

package maplimit

import "sync"

func MapLimit(items []int, limit int, fn func(int) (int, error)) ([]int, error) {
	results := make([]int, len(items))
	errs := make([]error, len(items))

	sem := make(chan struct{}, limit)
	var wg sync.WaitGroup

	for i, item := range items {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			results[i], errs[i] = fn(item)
		}()
	}

	wg.Wait()

	for _, err := range errs {
		if err != nil {
			return nil, err
		}
	}

	return results, nil
}
