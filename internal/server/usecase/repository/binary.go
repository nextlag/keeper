package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/nextlag/keeper/internal/entity"
	"github.com/nextlag/keeper/internal/server/usecase/repository/models"
	"github.com/nextlag/keeper/pkg/logger/l"
)

var errWrongBinaryOwner = errors.New("wrong binary owner or not found")

// GetBinaries retrieves all binaries associated with the specified user.
// Returns a slice of binaries and an error if something went wrong.
func (r *Repo) GetBinaries(ctx context.Context, user entity.User) (binaries []entity.Binary, err error) {
	var binariesFromDB []models.Binary

	if err = r.db.WithContext(ctx).
		Model(&models.Binary{}).
		Preload("Meta").
		Find(&binariesFromDB, "user_id = ?", user.ID).Error; err != nil {
		return nil, l.WrapErr(err)
	}

	if len(binariesFromDB) == 0 {
		return nil, nil
	}

	binaries = make([]entity.Binary, len(binariesFromDB))

	for index := range binariesFromDB {
		binaries[index].ID = binariesFromDB[index].ID
		binaries[index].Name = binariesFromDB[index].Name
		binaries[index].FileName = binariesFromDB[index].FileName
		for metaIndex := range binariesFromDB[index].Meta {
			binaries[index].Meta = append(binaries[index].Meta, entity.Meta{
				ID:    binariesFromDB[index].Meta[metaIndex].ID,
				Name:  binariesFromDB[index].Meta[metaIndex].Name,
				Value: binariesFromDB[index].Meta[metaIndex].Value,
			})
		}
	}

	return
}

// AddBinary inserts a new binary record into the database.
// Sets the ID of the binary after successful insertion.
func (r *Repo) AddBinary(ctx context.Context, binary *entity.Binary, userID uuid.UUID) error {
	newBinaryToDB := models.Binary{
		Name:     binary.Name,
		FileName: binary.FileName,
		UserID:   userID,
	}

	if err := r.db.WithContext(ctx).Create(&newBinaryToDB).Error; err != nil {
		return l.WrapErr(err)
	}
	binary.ID = newBinaryToDB.ID

	return nil
}

// GetBinary retrieves a single binary by its ID and ensures it belongs to the specified user.
// Returns the binary and an error if the binary is not found or if the user does not own it.
func (r *Repo) GetBinary(ctx context.Context, binaryID, userID uuid.UUID) (binary *entity.Binary, err error) {
	var binaryFromDB models.Binary
	if err = r.db.WithContext(ctx).
		Model(&models.Binary{}).
		Preload("Meta").
		Find(&binaryFromDB, binaryID).Error; err != nil {
		return nil, l.WrapErr(err)
	}

	if binaryFromDB.UserID != userID {
		err = errWrongBinaryOwner
		return nil, l.WrapErr(err)
	}

	var meta []entity.Meta
	if len(binaryFromDB.Meta) > 0 {
		meta = make([]entity.Meta, len(binaryFromDB.Meta))
		for index := range binaryFromDB.Meta {
			meta[index].ID = binaryFromDB.Meta[index].ID
			meta[index].Name = binaryFromDB.Meta[index].Name
			meta[index].Value = binaryFromDB.Meta[index].Value
		}
	}

	return &entity.Binary{
		ID:       binaryFromDB.ID,
		FileName: binaryFromDB.FileName,
		Meta:     meta,
	}, nil
}

// DelUserBinary deletes a binary record by its UUID if it belongs to the current user.
// Returns an error if the binary does not belong to the user or if deletion fails.
func (r *Repo) DelUserBinary(ctx context.Context, currentUser *entity.User, binaryUUID uuid.UUID) (err error) {
	var binaryFromDB models.Binary
	r.db.WithContext(ctx).Find(&binaryFromDB, binaryUUID)
	if binaryFromDB.UserID != currentUser.ID {
		err = errWrongBinaryOwner
		return l.WrapErr(err)
	}
	return l.WrapErr(r.db.Delete(&binaryFromDB).Error)
}

// AddBinaryMeta adds metadata to a binary record.
// Retrieves the updated binary after saving the metadata.
func (r *Repo) AddBinaryMeta(
	ctx context.Context,
	currentUser *entity.User,
	binaryUUID uuid.UUID,
	meta []entity.Meta,
) (*entity.Binary, error) {
	metaForDB := make([]models.MetaBinary, len(meta))
	for index := range meta {
		metaForDB[index].BinaryID = binaryUUID
		metaForDB[index].Name = meta[index].Name
		metaForDB[index].Value = meta[index].Value
	}

	if err := r.db.WithContext(ctx).Save(&metaForDB).Error; err != nil {
		return nil, l.WrapErr(err)
	}

	return r.GetBinary(ctx, binaryUUID, currentUser.ID)
}
