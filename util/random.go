package util

import (
	"math/rand"
	"strings"
	"time"
)

var alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// random generate a integer between min to max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min-1) // min + [0, max-min] -> [min, max]
}

// random generate a string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i <= n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomOwnerName() string {
	return RandomString(6)
}

func RandomBalance() int64 {
	return RandomInt(0, 1000)
}

func RandomCurrency() string {
	currency := []string{"USD", "EUR", "TWD"}
	n := len(currency)

	return currency[rand.Intn(n)]
}
