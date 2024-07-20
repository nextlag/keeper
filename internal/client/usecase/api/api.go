package api

type ClientAPI struct {
	serverURL string
}

func New(serverURL string) *ClientAPI {
	return &ClientAPI{
		serverURL: serverURL,
	}
}
