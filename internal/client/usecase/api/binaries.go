package api

import (
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"

	"github.com/nextlag/keeper/internal/entity"
)

const binaryEndpoint = "api/v1/user/binary"

func (api *ClientAPI) GetBinaries(accessToken string) (binaries []entity.Binary, err error) {
	if err = api.getEntities(&binaries, accessToken, binaryEndpoint); err != nil {
		return nil, err
	}
	return binaries, nil
}

func (api *ClientAPI) AddBinary(accessToken string, binary *entity.Binary, tmpFilePath string) error {
	var responseBinary entity.Binary
	client := resty.New()
	client.SetAuthToken(accessToken)
	file, err := os.Open(tmpFilePath)
	if err != nil {
		return fmt.Errorf("ClientAPI - AddBinary - %w ", err)
	}
	resp, err := client.R().
		SetHeader("Content-Type", "multipart/form-data").
		SetQueryParam("name", binary.Name).
		SetFileReader("file", binary.FileName, file).
		SetResult(&responseBinary).
		Post(fmt.Sprintf("%s/%s", api.serverURL, binaryEndpoint))
	if err != nil {
		return fmt.Errorf("ClientAPI - AddBinary - %w ", err)
	}

	if err = api.checkResCode(resp); err != nil {
		return errServer
	}
	binary.ID = responseBinary.ID
	if binary.Meta != nil {
		metaEndpoint := fmt.Sprintf("%s/%s/meta", binaryEndpoint, binary.ID.String())
		if err = api.addEntity(&binary.Meta, accessToken, metaEndpoint); err != nil {
			return fmt.Errorf("ClientAPI - AddBinary - AddMeta %w ", err)
		}
	}
	return nil
}

func (api *ClientAPI) DelBinary(accessToken, binaryID string) error {
	return api.delEntity(accessToken, binaryEndpoint, binaryID)
}

func (api *ClientAPI) DownloadBinary(accessToken, outputFilePath string, binary *entity.Binary) error {
	client := resty.New()
	client.SetAuthToken(accessToken)
	resp, err := client.R().
		SetOutput(outputFilePath).
		Get(fmt.Sprintf("%s/%s/%s", api.serverURL, binaryEndpoint, binary.ID.String()))
	if err != nil {
		return err
	}

	if err = api.checkResCode(resp); err != nil {
		return err
	}
	return nil
}
