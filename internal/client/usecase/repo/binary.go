package repo

import (
	"github.com/nextlag/keeper/internal/client/usecase/repo/models"
	"github.com/nextlag/keeper/internal/client/usecase/viewsets"
)

func (r *Repo) LoadBinaries() []viewsets.BinaryForList {
	userID := r.getUserID()
	var binaries []models.Binary
	r.db.
		Model(&models.Binary{}).
		Preload("Meta").
		Where("user_id", userID).Find(&binaries)
	if len(binaries) == 0 {
		return nil
	}

	binariesViewSet := make([]viewsets.BinaryForList, len(binaries))

	for index := range binaries {
		binariesViewSet[index].ID = binaries[index].ID
		binariesViewSet[index].Name = binaries[index].Name
		binariesViewSet[index].FileName = binaries[index].FileName
	}

	return binariesViewSet
}
