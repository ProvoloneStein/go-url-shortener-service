package repositories

import (
	"errors"
	"fmt"
	"math/rand"
)

var ErrUniqueViolation = errors.New("UniqueViolationError")
var ErrURLNotFound = errors.New("URLNotFound")
var ErrDeleted = errors.New("DeletedUrl")
var ErrShortURLExists = errors.New("ShortURLExists")
var ErrUserExists = errors.New("UserExists")

const (
	queryErrorMessage = "ошибка при обращении к бд"
	defaultRepoError  = "repository:"
	randomStringSize  = 10
	txRollbackError   = "transaction rollback error"
)

func defaultRepoErrWrapper(err error) error {
	return fmt.Errorf("%s %w", defaultRepoError, err)
}

func errWithVal(err error, val string) error {
	return fmt.Errorf("%w %s", err, val)
}

func RandomString() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, randomStringSize)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
