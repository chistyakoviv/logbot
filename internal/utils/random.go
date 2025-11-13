package utils

import (
	"crypto/rand"
	"math/big"
)

func RandFromSlice[T any](slice []T) T {
	max := big.NewInt(int64(len(slice)))
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		n = big.NewInt(0)
	}
	return slice[n.Int64()]
}
