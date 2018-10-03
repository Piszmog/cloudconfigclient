package cloudconfigclient

import (
	"github.com/Piszmog/cloudconfigclient/net"
	"github.com/pkg/errors"
	"net/http"
)

// ConfigClient contains the clients of the Config Servers.
type ConfigClient struct {
	Clients []CloudClient
}

// CloudClient interacts with the Config Server's REST APIs
type CloudClient interface {
	Get(uriVariables ...string) (resp *http.Response, err error)
}

// Client that wraps http.Client and the base Uri of the http client
type Client struct {
	configUri  string
	httpClient *http.Client
}

// Get performs a REST GET
func (client Client) Get(uriVariables ...string) (resp *http.Response, err error) {
	fullUrl := net.CreateUrl(client.configUri, uriVariables...)
	response, err := client.httpClient.Get(fullUrl)
	return response, errors.Wrapf(err, "failed to retrieve from %s", fullUrl)
}
