//go:build !solution

package dedup

// Deduplicate удаляет повторы, сохраняя порядок первых появлений.
func Deduplicate(input []string) []string {
	if input == nil {
		return nil
	}
	seen := make(map[string]struct{})
	result := make([]string, 0, len(input))
	for _, s := range input {
		if _, ok := seen[s]; !ok {
			seen[s] = struct{}{}
			result = append(result, s)
		}
	}
	return result
}
