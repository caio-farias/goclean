package main

func Filter[T any](arr []T, filterFn func(T) bool) []T {
	var filtered []T
	for _, item := range arr {
		if filterFn(item) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func FindOne[T any](arr []T, filterFn func(T) bool) T {
	subArr := Filter(arr, filterFn)
	return subArr[0]
}
