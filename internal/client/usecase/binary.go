package usecase

import (
	"os"

	"github.com/fatih/color"
	"github.com/google/uuid"

	"github.com/nextlag/keeper/internal/entity"
	"github.com/nextlag/keeper/internal/utils"
)

// AddBinary adds a binary file.
func (uc *ClientUseCase) AddBinary(userPassword string, binary *entity.Binary) {
	accessToken, err := uc.authorisationCheck(userPassword)
	if err != nil {
		color.Red("Authorization failed: %v", err)
		return
	}

	file, err := os.Stat(binary.FileName)
	if err != nil {
		color.Red("Error statting file %s: %v", binary.FileName, err)
		return
	}
	tmpFilePath := uc.cfg.FilesStorage.ClientLocation + file.Name()

	if err = utils.EncryptFile(userPassword, binary.FileName, tmpFilePath); err != nil {
		color.Red("Error encrypting file %s: %v", binary.FileName, err)
		return
	}

	if err = uc.clientAPI.AddBinary(accessToken, binary, tmpFilePath); err != nil {
		color.Red("Error adding binary file %s: %v", binary.FileName, err)
		return
	}

	if err = uc.repo.AddBinary(binary); err != nil {
		color.Red("Error saving binary file %s to repository: %v", binary.FileName, err)
		return
	}
	defer func() {
		if removeErr := os.Remove(tmpFilePath); removeErr != nil {
			color.Red("Failed to remove temporary file %s: %v", tmpFilePath, removeErr)
		}
	}()
	color.Green("Binary %v - %s saved successfully", binary.ID, binary.FileName)
}

// DelBinary deletes a binary file.
func (uc *ClientUseCase) DelBinary(userPassword, binaryID string) {
	accessToken, err := uc.authorisationCheck(userPassword)
	if err != nil {
		color.Red("Authorization failed: %v", err)
		return
	}

	binaryUUID, err := uuid.Parse(binaryID)
	if err != nil {
		color.Red("Error parsing binary ID %s: %v", binaryID, err)
		return
	}

	if err = uc.repo.DelBinary(binaryUUID); err != nil {
		color.Red("Error deleting binary file with ID %s from repository: %v", binaryID, err)
		return
	}

	if err = uc.clientAPI.DelBinary(accessToken, binaryID); err != nil {
		color.Red("Error deleting binary file %s with ID %s through API: %v", binaryID, binaryID, err)
		return
	}

	color.Green("Binary %q removed successfully", binaryID)
}

// GetBinary downloads and decrypts a binary file.
func (uc *ClientUseCase) GetBinary(userPassword, binaryID, filePath string) {
	accessToken, err := uc.authorisationCheck(userPassword)
	if err != nil {
		color.Red("Authorization failed: %v", err)
		return
	}

	binaryUUID, err := uuid.Parse(binaryID)
	if err != nil {
		color.Red("Error parsing binary ID %s: %v", binaryID, err)
		return
	}

	binary, err := uc.repo.GetBinaryByID(binaryUUID)
	if err != nil {
		color.Red("Error fetching binary file with ID %s from repository: %v", binaryID, err)
		return
	}

	tmpFilePath := filePath + binary.FileName

	defer func() {
		if removeErr := os.Remove(tmpFilePath); removeErr != nil {
			color.Red("Failed to remove temporary file %s: %v", tmpFilePath, removeErr)
		}
	}()

	if err = uc.clientAPI.DownloadBinary(accessToken, tmpFilePath, &binary); err != nil {
		color.Red("Error downloading binary file %s: %v", binary.FileName, err)
		return
	}

	if err = utils.DecryptFile(userPassword, tmpFilePath, filePath); err != nil {
		color.Red("Error decrypting file %s: %v", tmpFilePath, err)
		return
	}

	color.Green("File decrypted to %s", filePath)
	color.Green("File '%v' successfully downloaded", binary.FileName)
}

// loadBinaries loads binaries and saves them to the repository.
func (uc *ClientUseCase) loadBinaries(accessToken string) {
	binaries, err := uc.clientAPI.GetBinaries(accessToken)
	if err != nil {
		color.Red("Error fetching binaries with access token %s: %v", accessToken, err)
		return
	}

	if err = uc.repo.SaveBinaries(binaries); err != nil {
		color.Red("Error saving binaries to repository: %v", err)
		return
	}

	color.Green("Loaded %v binaries successfully", len(binaries))
}
