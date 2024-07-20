package api

import (
	"github.com/nextlag/keeper/internal/entity"
)

const cardsEndpoint = "api/v1/user/cards"

func (api *ClientAPI) GetCards(accessToken string) (cards []entity.Card, err error) {
	if err := api.getEntities(&cards, accessToken, cardsEndpoint); err != nil {
		return nil, err
	}

	return cards, nil
}

func (api *ClientAPI) AddCard(accessToken string, card *entity.Card) error {
	return api.addEntity(card, accessToken, cardsEndpoint)
}

func (api *ClientAPI) DelCard(accessToken, cardID string) error {
	return api.delEntity(accessToken, cardsEndpoint, cardID)
}
