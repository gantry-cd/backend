package utils

import "reflect"

// removeDuplocate removes duplicate elements from a slice
func RemoveDuplocate[T any](s []T) []T {
	for i := 0; i < len(s); i++ {
		for j := i + 1; j < len(s); j++ {
			if reflect.DeepEqual(s[i], s[j]) {
				s = append(s[:j], s[j+1:]...)
				j--
			}
		}
	}
	return s
}

// toPtr converts a value to a pointer
func ToPtr[T any](v T) *T {
	t := reflect.TypeOf(v)

	// ゼロ値の場合はnilを返す
	if reflect.DeepEqual(v, reflect.Zero(t).Interface()) {
		return nil
	}

	return &v
}
