package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MetaLogin represents metadata associated with a Login entity in the database.
type MetaLogin struct {
	gorm.Model
	ID      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name    string    // Name of the metadata
	Value   string    // Value associated with the metadata
	LoginID uuid.UUID // Foreign key reference to Login ID
}

// Login represents a user login entity in the database.
type Login struct {
	gorm.Model
	ID       uuid.UUID   `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name     string      `gorm:"size:100"` // Name of the login entry, limited to 100 characters
	URI      string      // URI associated with the login
	Login    string      // Login username
	Password string      // Login password
	UserID   uuid.UUID   // Foreign key reference to User ID
	Meta     []MetaLogin `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // Metadata associated with the login
}
