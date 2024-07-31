package entity

import "github.com/google/uuid"

// Binary represents a file.
type Binary struct {
	ID       uuid.UUID `json:"uuid"`      // Unique identifier.
	Name     string    `json:"name"`      // File name.
	FileName string    `json:"file_name"` // Filesystem name.
	Meta     []Meta    `json:"meta"`      // Associated metadata.
}
