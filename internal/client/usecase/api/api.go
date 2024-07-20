package api

type ClientAPI struct {
	host string
}

func New(host string) *ClientAPI {
	return &ClientAPI{
		host: host,
	}
}
