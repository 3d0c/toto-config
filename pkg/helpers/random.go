package helpers

import (
	"crypto/rand"
	"encoding/hex"
	"math/big"
)

// RandomSeed generates random integer in ragne
func RandomSeed(min, max int) int {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)-int64(min)+1))
	if err != nil {
		return max
	}
	return min + int(n.Int64())
}

// RandomString is used for tests only. That is why it can panic
func RandomString(n int) string {
	buf := make([]byte, n)

	if _, err := rand.Read(buf); err != nil {
		panic("Error generating random value " + err.Error())
	}

	return hex.EncodeToString(buf)
}
