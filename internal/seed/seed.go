// Package seed implements functionality for working with Ocarina of Time
// Randomizer seeds.
package seed

import (
	"crypto/rand"
	"math/big"
)

// GenerateSeedString generates a random string of a given length
// by reading from a cryptographically secure source of randomness.
//
// The generated string can be used to seed the Ocarina of time randomizer.
func GenerateSeedString(n int) (string, error) {
	const (
		letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		max     = int64(len(letters))
		bits    = 6
		mask    = 1<<bits - 1
	)

	b := make([]byte, n)
	for i := 0; i < n; {
		num, err := rand.Int(rand.Reader, big.NewInt(max))
		if err != nil {
			return "", err
		}
		if idx := int(num.Int64() & mask); idx < len(letters) {
			b[i] = letters[idx]
			i++
		}
	}
	return string(b), nil
}
