package v1

import "encoding/json"

func jsonError(err error) string {
	response, _ := json.Marshal(map[string]string{"error": err.Error()})
	return string(response)
}

func jsonResponse(message string) string {
	response, _ := json.Marshal(map[string]string{"status": message})
	return string(response)
}
