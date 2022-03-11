package security

import (
	"crypto/md5"
	"db-server/models"
	"fmt"
	"math/rand"
)

func ValidatePassword(password string, user models.User) bool {
	return user.PasswordHash == HashPassword(password)
}

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
