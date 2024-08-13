package usecase

import (
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/google/uuid"

	"github.com/nextlag/keeper/internal/entity"
	"github.com/nextlag/keeper/internal/utils"
)

func (uc *ClientUseCase) AddBinary(userPassword string, binary *entity.Binary) {
	accessToken, err := uc.authorisationCheck(userPassword)
	if err != nil {
		log.Printf("ClientUseCase - authorisationCheck - %v", err)
		return
	}
	file, err := os.Stat(binary.FileName)
	if err != nil {
		log.Printf("ClientUseCase - AddBinary - %v", err)
		return
	}
	tmpFilePath := uc.cfg.FilesStorage.ClientLocation + file.Name()

	if err = utils.EncryptFile(userPassword, binary.FileName, tmpFilePath); err != nil {
		log.Printf("ClientUseCase - EncryptFile - %v", err)
		return
	}

	if err = uc.clientAPI.AddBinary(accessToken, binary, tmpFilePath); err != nil {
		log.Printf("ClientUseCase - clientAPI.AddBinary - %v", err)
		return
	}

	if err = uc.repo.AddBinary(binary); err != nil {
		log.Printf("ClientUseCase - repo.AddBinary - %v", err)
		return
	}
	defer func() {
		if err = os.Remove(tmpFilePath); err != nil {
			log.Printf("ClientUseCase -  os.Remove - %v", err)
			return
		}
	}()
	color.Green("Binary %v - %s saved", binary.ID, binary.FileName)
}

func (uc *ClientUseCase) DelBinary(userPassword, binaryID string) {
	accessToken, err := uc.authorisationCheck(userPassword)
	if err != nil {
		return
	}
	binaryUUID, err := uuid.Parse(binaryID)
	if err != nil {
		color.Red(err.Error())
		log.Printf("ClientUseCase - DelBinary - uuid.Parse - %v", err)
		return
	}

	if err = uc.repo.DelBinary(binaryUUID); err != nil {
		log.Printf("ClientUseCase - repo.DelNote - %v", err)
		return
	}

	if err = uc.clientAPI.DelBinary(accessToken, binaryID); err != nil {
		log.Printf("ClientUseCase - clientAPI.DelBinary - %v", err)
		return
	}
	color.Green("Binary %q removed", binaryID)
}

func (uc *ClientUseCase) GetBinary(userPassword, binaryID, filePath string) {
	accessToken, err := uc.authorisationCheck(userPassword)
	if err != nil {
		return
	}

	binaryUUID, err := uuid.Parse(binaryID)
	if err != nil {
		log.Printf("ClientUseCase - GetBinary - uuid.Parse - %v", err)
		return
	}

	binary, err := uc.repo.GetBinaryByID(binaryUUID)
	if err != nil {
		log.Printf("ClientUseCase - GetBinary - GetBinaryByID - %v", err)
		return
	}

	tmpFilePath := filePath + binary.FileName

	defer func() {
		if removeErr := os.Remove(tmpFilePath); removeErr != nil {
			log.Printf("ClientUseCase - GetBinary - Failed to remove temporary file: %v", removeErr)
		}
	}()

	if err = uc.clientAPI.DownloadBinary(accessToken, tmpFilePath, &binary); err != nil {
		log.Printf("ClientUseCase - GetBinary - clientAPI.DownloadBinary - %v", err)
		return
	}

	if err = utils.DecryptFile(userPassword, tmpFilePath, filePath); err != nil {
		log.Printf("ClientUseCase - GetBinary - EncryptFile - %v", err)
		return
	}

	color.Green("File decrypted to %s", filePath)
	color.Green("File '%v' successfully downloaded", binary.FileName)
}

func (uc *ClientUseCase) loadBinaries(accessToken string) {
	binaries, err := uc.clientAPI.GetBinaries(accessToken)
	if err != nil {
		color.Red("Connection error: %v", err)
		return
	}

	if err = uc.repo.SaveBinaries(binaries); err != nil {
		log.Println(err)
		return
	}
	color.Green("Loaded %v binaries", len(binaries))
}
