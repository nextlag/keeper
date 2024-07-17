package entity

import "github.com/google/uuid"

type (
	Meta struct {
		ID    uuid.UUID `json:"uuid"`
		Name  string    `json:"name"`
		Value string    `json:"value"`
	}
)
