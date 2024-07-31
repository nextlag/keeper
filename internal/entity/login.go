package entity

import "github.com/google/uuid"

// Login represents a user login entry with associated metadata.
type Login struct {
	ID       uuid.UUID `json:"uuid"`     // Unique identifier.
	Name     string    `json:"name"`     // Name of the login entry.
	Login    string    `json:"login"`    // Login username or identifier.
	Password string    `json:"password"` // Password for the login.
	URI      string    `json:"uri"`      // URI or website related to the login.
	Meta     []Meta    `json:"meta"`     // Associated metadata.
}
