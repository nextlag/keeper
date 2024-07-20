package entity

import "github.com/google/uuid"

type Binary struct {
	ID       uuid.UUID `json:"uuid"`
	Name     string    `json:"name"`
	FileName string    `json:"file_name"`
	Meta     []Meta    `json:"meta"`
}
