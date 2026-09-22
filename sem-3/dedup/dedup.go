//go:build !solution

package dedup

// Deduplicate удаляет повторы, сохраняя порядок первых появлений.
func Deduplicate(input []string) []string {
	if input == nil {
		return nil
	}

	seen := make(map[string]struct{})
	result := make([]string, 0)

	for _, str := range input {
		if _, exists := seen[str]; !exists {
			seen[str] = struct{}{}
			result = append(result, str)
		}
	}

	return result
}
