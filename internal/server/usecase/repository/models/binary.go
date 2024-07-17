package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MetaBinary represents metadata associated with a Binary entity.
type MetaBinary struct {
	gorm.Model
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name     string    // Name of the metadata
	Value    string    // Value associated with the metadata
	BinaryID uuid.UUID `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // Foreign key reference to Binary ID
}

// Binary represents a binary entity stored in the database.
type Binary struct {
	gorm.Model
	ID       uuid.UUID    `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name     string       // Name of the binary data
	FileName string       // File name associated with the binary data
	UserID   uuid.UUID    // Foreign key reference to User ID
	Meta     []MetaBinary // Metadata associated with the binary data
}
