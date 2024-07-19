package entity

import "github.com/google/uuid"

type SecretNote struct {
	ID   uuid.UUID `json:"uuid"`
	Name string    `json:"name"`
	Note string    `json:"note"`
	Meta []Meta    `json:"meta"`
}
