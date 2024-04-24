package utils

func Filter[T any](items []T, test func(T) bool) (results []T) {
	for _, s := range items {
		if test(s) {
			results = append(results, s)
		}
	}
	return results
}
