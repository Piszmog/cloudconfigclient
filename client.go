package cloudconfigclient

import (
	"encoding/json"
	"fmt"
	"github.com/Piszmog/cloudconfigclient/net"
	"io/ioutil"
	"net/http"
)

// NotFoundError is a special error that is used to propagate 404s.
type NotFoundError struct {
}

// Error return the error message.
func (r NotFoundError) Error() string {
	return "failed to find resource"
}

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
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve from %s: %w", fullUrl, err)
	}
	return response, nil
}

func getResource(client CloudClient, dest interface{}, uriVariables ...string) error {
	resp, err := client.Get(uriVariables...)
	if err != nil {
		return err
	}
	defer closeResource(resp.Body)
	if resp.StatusCode == 404 {
		return &NotFoundError{}
	}
	if resp.StatusCode != 200 {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read body with status code %d: %w", resp.StatusCode, err)
		}
		return fmt.Errorf("server responded with status code %d and body %s", resp.StatusCode, b)
	}
	if err := json.NewDecoder(resp.Body).Decode(dest); err != nil {
		return fmt.Errorf("failed to decode response from url: %w", err)
	}
	return nil
}
