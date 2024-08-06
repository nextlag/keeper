package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MetaNote represents a metadata associated with a Note.
type MetaNote struct {
	gorm.Model
	ID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name   string    // Name of the metadata
	Value  string    // Value associated with the metadata
	NoteID uuid.UUID // Foreign key reference to Note ID
}

// Note represents a note entity with associated metadata.
type Note struct {
	gorm.Model
	ID     uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name   string     `gorm:"size:100"` // Name of the note, limited to 100 characters
	Note   string     // Content of the note
	UserID uuid.UUID  // Foreign key reference to User ID
	Meta   []MetaNote `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // Metadata associated with the note
}
