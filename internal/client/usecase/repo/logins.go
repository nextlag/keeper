package repo

import (
	"github.com/nextlag/keeper/internal/client/usecase/repo/models"
	"github.com/nextlag/keeper/internal/client/usecase/viewsets"
)

func (r *Repo) LoadLogins() []viewsets.LoginForList {
	userID := r.getUserID()
	var logins []models.Login
	r.db.
		Model(&models.Login{}).
		Preload("Meta").
		Where("user_id", userID).Find(&logins)
	if len(logins) == 0 {
		return nil
	}

	loginsViewSet := make([]viewsets.LoginForList, len(logins))

	for index := range logins {
		loginsViewSet[index].ID = logins[index].ID
		loginsViewSet[index].Name = logins[index].Name
		loginsViewSet[index].URI = logins[index].URI
	}
	return loginsViewSet
}
