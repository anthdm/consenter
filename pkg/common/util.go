package common

import (
	"crypto/sha256"
	"math/rand"
	"time"
)

// RandInt return a random number between min and max.
func RandInt(min, max int) int {
	return min + rand.Intn(max-min)
}

// Hash256 computes the double sha256 hash of the given bytes.
func Hash256(b []byte) []byte {
	sha := sha256.New()
	sha.Write(b)
	hash := sha.Sum(nil)
	sha.Write(hash)
	return sha.Sum(nil)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
