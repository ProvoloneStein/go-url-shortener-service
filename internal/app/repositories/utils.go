package repositories

import (
	"errors"
	"fmt"
	"math/rand"
)

var ErrorUniqueViolation = errors.New("UniqueViolationError")
var UrlNotFound = errors.New("UrlNotFound")

type ValueError struct {
	Value string
	Err   error
}

// Error добавляет поддержку интерфейса error для типа ValueError.
func (ve *ValueError) Error() string {
	return fmt.Sprintf("%v - %v", ve.Value, ve.Err)
}

func NewValueError(value string, err error) error {
	return &ValueError{
		Value: value,
		Err:   err,
	}
}

func randomString() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, 10)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
