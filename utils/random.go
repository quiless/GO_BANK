package utils

import (
	"math/rand"
	"strings"
	"time"
)

const (
	alphabet = "abcdefghijklkmnopqrstuvwxyz"
)

var currencies = []string{"EU", "USD", "CAD"}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Generates a random integer between min and max
func GenerateInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// Generate a random string of length n
func GenerateString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte((c))
	}

	return sb.String()
}

// Random currency type
func GetCurrencyType() string {
	x := len(currencies)
	return currencies[(rand.Intn(x))]
}
