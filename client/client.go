package client

import (
	"github.com/Piszmog/cloudconfigclient/net"
	"net/http"
)

type CloudClient interface {
	GetFullUrl(uriVariables ...string) string
	Get(uriVariables ...string) (resp *http.Response, err error)
}

type ConfigClient struct {
	Clients []Client
}

type Client struct {
	configUri  string
	httpClient *http.Client
}

func (client *Client) GetFullUrl(uriVariables ...string) string {
	return net.CreateUrl(client.configUri, uriVariables...)
}

func (client *Client) Get(uriVariables ...string) (resp *http.Response, err error) {
	fullUrl := net.CreateUrl(client.configUri, uriVariables...)
	return client.httpClient.Get(fullUrl)
}
