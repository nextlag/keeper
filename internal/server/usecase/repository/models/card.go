package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MetaCard represents metadata associated with a Card entity in the database.
type MetaCard struct {
	gorm.Model
	ID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name   string    // Name of the metadata
	Value  string    // Value associated with the metadata
	CardID uuid.UUID // Foreign key reference to Card ID
}

// Card represents a credit/debit card entity in the database.
type Card struct {
	gorm.Model
	ID              uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name            string     `gorm:"size:100"` // Name of the card, limited to 100 characters
	CardHolderName  string     // Name of the cardholder
	Number          string     // Card number
	Brand           string     // Brand of the card (e.g., Visa, Mastercard)
	ExpirationMonth string     // Expiration month of the card (MM format)
	ExpirationYear  string     // Expiration year of the card (YYYY format)
	SecurityCode    string     // Security code (CVV) of the card
	UserID          uuid.UUID  // Foreign key reference to User ID
	Meta            []MetaCard `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // Metadata associated with the card
}
