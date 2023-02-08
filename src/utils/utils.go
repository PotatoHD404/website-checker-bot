package utils

import "reflect"

// type with == operator

func Contains[T any](els []T, el1 T) bool {
	for _, el2 := range els {
		if reflect.DeepEqual(el1, el2) {
			return true
		}
	}
	return false
}
