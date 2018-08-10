package client

import (
	"github.com/Piszmog/cloudconfigclient/net"
	"github.com/pkg/errors"
	"net/http"
)

type ConfigClient struct {
	Clients []Client
}

type CloudClient interface {
	Get(uriVariables ...string) (resp *http.Response, err error)
}

type Client struct {
	configUri  string
	httpClient *http.Client
}

func (client *Client) Get(uriVariables ...string) (resp *http.Response, err error) {
	fullUrl := net.CreateUrl(client.configUri, uriVariables...)
	response, err := client.httpClient.Get(fullUrl)
	return response, errors.Wrapf(err, "failed to retrieve from %s", fullUrl)
}
