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

type GormErr struct {
	Code    string `json:"Code"`
	Message string `json:"Message"`
}

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

func ParseServerError(body []byte) string {
	var errMessage struct {
		Message string `json:"error"`
	}

	// Check if the body is compressed with gzip
	if isGzipped(body) {
		unzippedBody, err := unzip(body)
		if err != nil {
			log.Printf("Failed to unzip response body: %v", err)
			return ""
		}
		body = unzippedBody
	}

	// Attempt to unmarshal JSON
	if err := json.Unmarshal(body, &errMessage); err != nil {
		log.Printf("Failed to parse error message as JSON: %v", err)
		return string(body)
	}

	if errMessage.Message == "" {
		log.Println("Error message is empty")
	}

	return errMessage.Message
}

// isGzipped checks if the body is gzipped.
func isGzipped(body []byte) bool {
	return len(body) > 2 && body[0] == 0x1f && body[1] == 0x8b
}

// unzip decompresses gzipped data.
func unzip(body []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return io.ReadAll(reader)
}
