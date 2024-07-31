package entity

import "github.com/google/uuid"

// SecretNote represents a note with associated metadata.
type SecretNote struct {
	ID   uuid.UUID `json:"uuid"` // Unique identifier for the note.
	Name string    `json:"name"` // Name or title of the note.
	Note string    `json:"note"` // Content of the note.
	Meta []Meta    `json:"meta"` // Associated metadata for the note.
}
