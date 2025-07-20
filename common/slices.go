package common

func FilterSlice[T any](slice []T, predicate func(T) bool) []T {
	var result []T
	for _, item := range slice {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

func FindInSlice[T any](slice []T, predicate func(T) bool) (T, bool) {
	var defaultVal T
	for _, v := range slice {
		if predicate(v) {
			return v, true
		}
	}
	return defaultVal, false
}
