package usecase

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"

	"github.com/google/uuid"

	"github.com/nextlag/keeper/internal/entity"
	"github.com/nextlag/keeper/internal/utils"
	"github.com/nextlag/keeper/pkg/logger/l"
)

// GetBinaries retrieves all binaries associated with the given user.
func (uc *UseCase) GetBinaries(ctx context.Context, user entity.User) ([]entity.Binary, error) {
	return uc.repo.GetBinaries(ctx, user)
}

// AddBinary adds a new binary file to the storage and database.
func (uc *UseCase) AddBinary(
	ctx context.Context,
	binary *entity.Binary,
	file *multipart.FileHeader,
	userID uuid.UUID,
) error {

	userDirectory := uc.cfg.FilesStorage.Location + "/" + userID.String()
	if err := uc.repo.AddBinary(ctx, binary, userID); err != nil {
		return l.WrapErr(err)
	}
	if err := utils.SaveUploadedFile(file, binary.ID.String(), userDirectory); err != nil {
		uc.log.Debug("error", l.ErrAttr(err))
		return l.WrapErr(err)
	}
	return nil
}

// GetUserBinary retrieves the file path for a binary owned by the given user.
func (uc *UseCase) GetUserBinary(
	ctx context.Context,
	currentUser *entity.User,
	binaryUUID uuid.UUID,
) (filePath string, err error) {

	binary, err := uc.repo.GetBinary(ctx, binaryUUID, currentUser.ID)
	if err != nil {
		return "", l.WrapErr(err)
	}
	return fmt.Sprintf(
		"%s/%s/%s",
		uc.cfg.FilesStorage.Location,
		currentUser.ID.String(),
		binary.ID), nil
}

// DelUserBinary deletes a binary file from the storage and database.
func (uc *UseCase) DelUserBinary(
	ctx context.Context,
	currentUser *entity.User,
	binaryUUID uuid.UUID,
) error {

	err := uc.repo.DelUserBinary(ctx, currentUser, binaryUUID)
	if err != nil {
		return l.WrapErr(err)
	}

	filePath := fmt.Sprintf(
		"%s/%s/%s",
		uc.cfg.FilesStorage.Location,
		currentUser.ID.String(),
		binaryUUID.String())

	return os.Remove(filePath)
}

// AddBinaryMeta adds metadata to an existing binary owned by the given user.
func (uc *UseCase) AddBinaryMeta(
	ctx context.Context,
	currentUser *entity.User,
	binaryUUID uuid.UUID,
	meta []entity.Meta,
) (*entity.Binary, error) {
	return uc.repo.AddBinaryMeta(
		ctx,
		currentUser,
		binaryUUID,
		meta,
	)
}
