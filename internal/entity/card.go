package entity

import "github.com/google/uuid"

// Card represents a payment card with details and metadata.
type Card struct {
	ID              uuid.UUID `json:"uuid"`             // Unique identifier.
	Name            string    `json:"name"`             // Card name.
	CardHolderName  string    `json:"card_holder_name"` // Cardholder's name.
	Number          string    `json:"number"`           // Card number.
	Brand           string    `json:"brand"`            // Card brand.
	ExpirationMonth string    `json:"expiration_month"` // Expiration month.
	ExpirationYear  string    `json:"expiration_year"`  // Expiration year.
	SecurityCode    string    `json:"security_code"`    // Security code (CVV).
	Meta            []Meta    `json:"meta"`             // Associated metadata.
}
