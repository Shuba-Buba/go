//go:build !solution

package dedup

// Deduplicate удаляет повторы, сохраняя порядок первых появлений.
func Deduplicate(input []string) []string {
	if input == nil {
		return nil
	} else if len(input) == 0 {
		return input
	}

	result := make([]string, 0, len(input))
	unique_words := make(map[string]struct{})
	for _, word := range input {
		if _, ok := unique_words[word]; ok {
			continue
		}
		unique_words[word] = struct{}{}
		result = append(result, word)
	}
	return result
}
