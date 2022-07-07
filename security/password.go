package security

import (
	"crypto/md5"
	"fmt"
	"math/rand"
)

func HashPassword(password string) string {

	return fmt.Sprintf("%x", md5.Sum([]byte(password)))
}

func GenerateRandomString(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
