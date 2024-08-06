package entity

import "github.com/google/uuid"

// User represents a user entity with authentication details.
type User struct {
	ID       uuid.UUID `json:"uuid"`  // Unique identifier for the user.
	Email    string    `json:"email"` // Email address of the user.
	Password string    `json:"-"`     // Password for the user (not serialized).
}
