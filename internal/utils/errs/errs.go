package errs

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"io"
	"log"
)

var (
	ErrWrongEmail           = errors.New("incorrect email given")
	ErrEmailAlreadyExists   = errors.New("given email already exists")
	ErrWrongCredentials     = errors.New("wrong credentials have been given")
	ErrTokenValidation      = errors.New("token validation error")
	ErrUnexpectedError      = errors.New("some unexpected error")
	ErrWrongOwnerOrNotFound = errors.New("wrong owner or not found")
)

// GormErr represents an error structure typically returned by GORM.
type GormErr struct {
	Code    string `json:"Code"`
	Message string `json:"Message"`
}

// ParsePostgresErr parses a PostgreSQL error into a GormErr structure.
// It attempts to convert the provided error into a JSON format and then
// unmarshals it into a GormErr object. If the error cannot be parsed or
// unmarshalled, an empty GormErr is returned.
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

// ParseServerError extracts an error message from a response body.
// It first checks if the body is gzipped and, if so, decompresses it.
// After ensuring the body is in its uncompressed form, it attempts
// to unmarshal the body into a structure with a "error" field.
// If unmarshalling fails or the error message is empty, the raw body
// content is returned.
func ParseServerError(body []byte) string {
	var errMessage struct {
		Message string `json:"error"`
	}

	if isGzipped(body) {
		unzippedBody, err := unzip(body)
		if err != nil {
			log.Printf("Failed to unzip response body: %v", err)
			return ""
		}
		body = unzippedBody
	}

	if err := json.Unmarshal(body, &errMessage); err != nil {
		log.Printf("Failed to parse error message as JSON: %v", err)
		return string(body)
	}

	if errMessage.Message == "" {
		log.Println("Error message is empty")
	}

	return errMessage.Message
}

// isGzipped determines if the provided byte slice is gzipped based
// on its initial bytes. It checks for the gzip magic number at the
// beginning of the byte slice.
func isGzipped(body []byte) bool {
	return len(body) > 2 && body[0] == 0x1f && body[1] == 0x8b
}

// unzip decompresses gzipped data into a plain byte slice.
// It uses a gzip reader to read from the provided gzipped data and
// returns the decompressed byte slice.
func unzip(body []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return io.ReadAll(reader)
}
