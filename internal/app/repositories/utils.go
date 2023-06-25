package repositories

import (
	"errors"
	"math/rand"
)

var ErrorUniqueViolation = errors.New("UniqueViolationError")
var ErrURLNotFound = errors.New("URLNotFound")
var ErrShortURLExists = errors.New("ShortURLExists")

func RandomString() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, 10)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
