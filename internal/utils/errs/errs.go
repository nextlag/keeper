package errs

import (
	"encoding/json"
	"errors"
)

// Predefined error variables for common error scenarios.
var (
	ErrWrongEmail         = errors.New("incorrect email given")
	ErrEmailAlreadyExists = errors.New("given email already exists")
	ErrWrongCredentials   = errors.New("wrong credentials have been given")
	ErrTokenValidation    = errors.New("token validation error")
	ErrUnexpectedError    = errors.New("some unexpected error")
)

// GormErr represents a structured error response from the database layer.
type GormErr struct {
	Code    string `json:"Code"`
	Message string `json:"Message"`
}

// ParsePostgresErr parses a PostgreSQL database error into a structured GormErr.
func ParsePostgresErr(dbErr error) (newError GormErr) {
	byteErr, err := json.Marshal(dbErr)
	if err != nil {
		return
	}

	if err = json.Unmarshal(byteErr, &newError); err != nil {
		return GormErr{}
	}
	return
}
