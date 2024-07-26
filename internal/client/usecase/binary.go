package usecase

import (
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/google/uuid"

	"github.com/nextlag/keeper/internal/entity"
	"github.com/nextlag/keeper/internal/utils"
)

const filePath = "tmp/"

func (uc *ClientUseCase) AddBinary(userPassword string, binary *entity.Binary) {
	accessToken, err := uc.authorisationCheck(userPassword)
	if err != nil {
		log.Fatalf("ClientUseCase - authorisationCheck - %v", err)
	}
	file, err := os.Stat(binary.FileName)
	if err != nil {
		log.Fatalf("ClientUseCase - AddBinary - %v", err)
	}
	tmpFilePath := filePath + file.Name()

	if err = utils.EncryptFile(userPassword, binary.FileName, tmpFilePath); err != nil {
		log.Fatalf("ClientUseCase - EncryptFile - %v", err)
	}

	if err = uc.clientAPI.AddBinary(accessToken, binary, tmpFilePath); err != nil {
		log.Fatalf("ClientUseCase - clientAPI.AddBinary - %v", err)
	}

	if err = uc.repo.AddBinary(binary); err != nil {
		log.Fatalf("ClientUseCase - repo.AddBinary - %v", err)
	}

	if err = os.Rename(
		tmpFilePath,
		uc.cfg.FilesStorage.Location+"/"+binary.ID.String()); err != nil {
		log.Fatalf("ClientUseCase - os.Rename - %v", err)
	}
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
		log.Fatalf("ClientUseCase - DelBinary - uuid.Parse - %v", err)
	}

	if err = uc.repo.DelBinary(binaryUUID); err != nil {
		log.Fatalf("ClientUseCase - repo.DelNote - %v", err)
	}

	if err = uc.clientAPI.DelBinary(accessToken, binaryID); err != nil {
		log.Fatalf("ClientUseCase - clientAPI.DelBinary - %v", err)
	}
	if err = os.Remove(uc.cfg.FilesStorage.Location + "/" + binaryID); err != nil {
		log.Fatalf("ClientUseCase -  os.Remove - %v", err)
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
		log.Fatalf("ClientUseCase - GetBinary - uuid.Parse - %v", err)
	}

	binary, err := uc.repo.GetBinaryByID(binaryUUID)
	if err != nil {
		log.Fatalf("ClientUseCase - GetBinary - GetBinaryByID - %v", err)
	}

	tmpFilePath := filePath + binary.FileName

	defer func() {
		if removeErr := os.Remove(tmpFilePath); removeErr != nil {
			log.Printf("ClientUseCase - GetBinary - Failed to remove temporary file: %v", removeErr)
		}
	}()

	if err = uc.clientAPI.DownloadBinary(accessToken, tmpFilePath, &binary); err != nil {
		log.Fatalf("ClientUseCase - GetBinary - clientAPI.DownloadBinary - %v", err)
	}

	if err = utils.DecryptFile(userPassword, tmpFilePath, filePath); err != nil {
		log.Fatalf("ClientUseCase - GetBinary - EncryptFile - %v", err)
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
