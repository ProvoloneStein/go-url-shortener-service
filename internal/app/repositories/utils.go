package repositories

import (
	"errors"
	"fmt"
	"math/rand"
)

var ErrorUniqueViolation = errors.New("UniqueViolationError")
var ErrURLNotFound = errors.New("URLNotFound")

type ValueError struct {
	Value string
	Err   error
}

// Error добавляет поддержку интерфейса error для типа ValueError.

func NewValueError(value string, err error) error {
	return &ValueError{
		Value: value,
		Err:   err,
	}
}

func (ve *ValueError) Error() string {
	return fmt.Sprintf("%v - %v", ve.Value, ve.Err)
}

func (ve *ValueError) Unwrap() error {
	return ve.Err
}

func randomString() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, 10)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
