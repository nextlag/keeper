package entity

import "github.com/google/uuid"

// Meta represents additional metadata associated with an entity.
type Meta struct {
	ID    uuid.UUID `json:"uuid"`  // Unique identifier for the metadata.
	Name  string    `json:"name"`  // Name or type of the metadata.
	Value string    `json:"value"` // Value of the metadata.
}
