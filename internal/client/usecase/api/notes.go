package api

import "github.com/nextlag/keeper/internal/entity"

const notesEndpoint = "api/v1/user/notes"

func (api *ClientAPI) GetNotes(accessToken string) (notes []entity.SecretNote, err error) {
	if err := api.getEntities(&notes, accessToken, notesEndpoint); err != nil {
		return nil, err
	}

	return notes, nil
}

func (api *ClientAPI) AddNote(accessToken string, note *entity.SecretNote) error {
	return api.addEntity(note, accessToken, notesEndpoint)
}

func (api *ClientAPI) DelNote(accessToken, noteID string) error {
	return api.delEntity(accessToken, notesEndpoint, noteID)
}
