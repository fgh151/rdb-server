package utils

import (
	"crypto/md5"
	"fmt"
	"math/rand"
)

func HashPassword(password string) string {

	return fmt.Sprintf("%x", md5.Sum([]byte(password)))
}

func GenerateRandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func ValidateKey(k1 string, k2 string) bool {
	return k1 == k2
}
