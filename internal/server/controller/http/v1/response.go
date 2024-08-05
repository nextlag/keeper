package v1

import (
	"encoding/json"
)

type response struct {
	Error  string `json:"error,omitempty" example:"message"`
	Status string `json:"status,omitempty" example:"message"`
}

func jsonError(err error) string {
	r := response{Error: err.Error()}
	jsonData, _ := json.Marshal(r)
	return string(jsonData)
}

func jsonResponse(message string) string {
	r := response{Status: message}
	jsonData, _ := json.Marshal(r)
	return string(jsonData)
}
