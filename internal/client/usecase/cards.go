package usecase

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/google/uuid"

	"github.com/nextlag/keeper/internal/entity"
	"github.com/nextlag/keeper/internal/utils"
)

// AddCard adds a new card for the user.
func (uc *ClientUseCase) AddCard(userPassword string, card *entity.Card) {
	accessToken, err := uc.authorisationCheck(userPassword)
	if err != nil {
		color.Red("Authorization failed for user with provided password: %v", err)
		return
	}
	uc.encryptCard(userPassword, card)

	if err = uc.clientAPI.AddCard(accessToken, card); err != nil {
		color.Red("Error adding card %q with access token %s: %v", card.Name, accessToken, err)
		return
	}

	if err = uc.repo.AddCard(card); err != nil {
		color.Red("Error adding card %q to repository: %v", card.Name, err)
		return
	}

	color.Green("Card %q added successfully, ID: %v", card.Name, card.ID)
}

// encryptCard encrypts card data using the user's password.
func (uc *ClientUseCase) encryptCard(userPassword string, card *entity.Card) {
	card.Number = utils.Encrypt(userPassword, card.Number)
	card.SecurityCode = utils.Encrypt(userPassword, card.SecurityCode)
	card.ExpirationMonth = utils.Encrypt(userPassword, card.ExpirationMonth)
	card.ExpirationYear = utils.Encrypt(userPassword, card.ExpirationYear)
	card.CardHolderName = utils.Encrypt(userPassword, card.CardHolderName)
}

// decryptCard decrypts card data using the user's password.
func (uc *ClientUseCase) decryptCard(userPassword string, card *entity.Card) {
	card.Number = utils.Decrypt(userPassword, card.Number)
	card.SecurityCode = utils.Decrypt(userPassword, card.SecurityCode)
	card.ExpirationMonth = utils.Decrypt(userPassword, card.ExpirationMonth)
	card.ExpirationYear = utils.Decrypt(userPassword, card.ExpirationYear)
	card.CardHolderName = utils.Decrypt(userPassword, card.CardHolderName)
}

// ShowCard displays the card by its ID.
func (uc *ClientUseCase) ShowCard(userPassword, cardID string) {
	if !uc.verifyPassword(userPassword) {
		color.Red("Password verification failed")
		return
	}

	cardUUID, err := uuid.Parse(cardID)
	if err != nil {
		color.Red("Error parsing card ID %s: %v", cardID, err)
		return
	}

	card, err := uc.repo.GetCardByID(cardUUID)
	if err != nil {
		color.Red("Error fetching card with ID %s from repository: %v", cardID, err)
		return
	}

	uc.decryptCard(userPassword, &card)
	yellow := color.New(color.FgYellow).SprintFunc()
	fmt.Printf("ID: %s\nName: %s\nCardHolderName: %s\nNumber: %s\nBrand: %s\nExpiration: %s/%s\nCode: %s\nMeta: %v\n",
		yellow(card.ID),
		yellow(card.Name),
		yellow(card.CardHolderName),
		yellow(card.Number),
		yellow(card.Brand),
		yellow(card.ExpirationMonth),
		yellow(card.ExpirationYear),
		yellow(card.SecurityCode),
		yellow(card.Meta),
	)
}

// DelCard deletes the card by its ID.
func (uc *ClientUseCase) DelCard(userPassword, cardID string) {
	accessToken, err := uc.authorisationCheck(userPassword)
	if err != nil {
		color.Red("Authorization failed for user with provided password: %v", err)
		return
	}

	cardUUID, err := uuid.Parse(cardID)
	if err != nil {
		color.Red("Error parsing card ID %s: %v", cardID, err)
		return
	}

	if err = uc.repo.DelCard(cardUUID); err != nil {
		color.Red("Error deleting card with ID %s from repository: %v", cardID, err)
		return
	}

	if err = uc.clientAPI.DelCard(accessToken, cardID); err != nil {
		color.Red("Error deleting card %s with access token %s: %v", cardID, accessToken, err)
		return
	}

	color.Green("Card %q removed successfully", cardID)
}

// loadCards loads cards using the API and saves them to the repository.
func (uc *ClientUseCase) loadCards(accessToken string) {
	cards, err := uc.clientAPI.GetCards(accessToken)
	if err != nil {
		color.Red("Error fetching cards with access token %s: %v", accessToken, err)
		return
	}

	if err = uc.repo.SaveCards(cards); err != nil {
		color.Red("Error saving cards to repository: %v", err)
		return
	}

	color.Green("Loaded %v cards successfully", len(cards))
}
