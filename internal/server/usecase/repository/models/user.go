package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// User represents a user entity in the database.
type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Email     string    `gorm:"unique;uniqueIndex;not null"` // Email address of the user
	Password  string    `gorm:"not null"`                    // Password hash of the user
	CreatedAt time.Time // Timestamp when the user was created
	UpdatedAt time.Time // Timestamp when the user was last updated
	Cards     []Card    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // List of cards associated with the user
	Logins    []Login   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // List of logins associated with the user
	Notes     []Note    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // List of notes associated with the user
	Binary    []Binary  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // List of binary data associated with the user
}

// ToString returns a formatted string representation of the user.
func (user *User) ToString() string {
	return fmt.Sprintf("id: %v\nemail: %s", user.ID, user.Email)
}
