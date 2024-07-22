package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
	"golang.org/x/exp/slices"

	"github.com/nextlag/keeper/internal/utils/errs"
)

type ClientAPI struct {
	serverURL string
}

func New(serverURL string) *ClientAPI {
	return &ClientAPI{
		serverURL: serverURL,
	}
}

// addEntity sends a POST request to add a new entity to the server.
func (api *ClientAPI) addEntity(entity any, accessToken, endpoint string) error {
	client := resty.New()
	client.SetAuthToken(accessToken)
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(entity).
		SetResult(entity).
		Post(fmt.Sprintf("%s/%s", api.serverURL, endpoint))
	if err != nil {
		log.Fatalf("ClientAPI - client.R - %v ", err)
	}
	if err := api.checkResCode(resp); err != nil {
		return errServer
	}

	return nil
}

// getEntities sends a GET request to retrieve entities from the server.
func (api *ClientAPI) getEntities(entity any, accessToken, endpoint string) error {
	client := resty.New()
	client.SetAuthToken(accessToken)
	resp, err := client.R().
		SetResult(entity).
		Get(fmt.Sprintf("%s/%s", api.serverURL, endpoint))
	if err != nil {
		log.Println(err)
		return err
	}

	if err := api.checkResCode(resp); err != nil {
		return err
	}

	return nil
}

// checkResCode checks the response code and returns an error if the status code indicates a failure.
func (api *ClientAPI) checkResCode(resp *resty.Response) error {
	badCodes := []int{http.StatusBadRequest, http.StatusInternalServerError, http.StatusUnauthorized}
	if slices.Contains(badCodes, resp.StatusCode()) {
		errMessage := errs.ParseServerError(resp.Body())
		color.Red("Server error: %s", errMessage)
		return errServer
	}

	return nil
}

// delEntity sends a DELETE request to delete an entity from the server.
func (api *ClientAPI) delEntity(accessToken, endpoint, id string) error {
	client := resty.New()
	client.SetAuthToken(accessToken)
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		Delete(fmt.Sprintf("%s/%s/%s", api.serverURL, endpoint, id))
	if err != nil {
		log.Fatalf("ClientAPI - client.R - %v ", err)
	}
	if err := api.checkResCode(resp); err != nil {
		return errServer
	}

	return nil
}
