package repositories

import (
	"errors"
)

var ErrorUniqueViolation = errors.New("UniqueViolationError")
var ErrURLNotFound = errors.New("URLNotFound")
var ErrShortURLExists = errors.New("ShortURLExists")
