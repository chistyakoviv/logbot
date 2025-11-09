package utils

import (
	"math/rand"
)

func RandFromSlice[T any](slice []T) T {
	return slice[rand.Intn(len(slice))]
}
