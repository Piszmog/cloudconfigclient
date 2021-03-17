package cloudconfigclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
)

// ConfigClient contains the clients of the Config Servers.
type ConfigClient struct {
	// Clients are all the config server clients
	Clients []CloudClient
}

// CloudClient interacts with the Config Server's REST APIs
type CloudClient interface {
	Get(uriPaths ...string) (*http.Response, error)
}

// Client that wraps http.Client and the base Uri of the http client
type Client struct {
	// ConfigUri is the uri of the config server
	ConfigUri string
	// HttpClient is the HTTP client to use to make the HTTP requests and handle responses
	HttpClient *http.Client
}

// Get performs a REST GET
func (c Client) Get(uriPaths ...string) (*http.Response, error) {
	fullUrl, err := createUrl(c.ConfigUri, uriPaths...)
	if err != nil {
		return nil, fmt.Errorf("failed to create url: %w", err)
	}
	response, err := c.HttpClient.Get(fullUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve from %s: %w", fullUrl, err)
	}
	return response, nil
}

func createUrl(baseUrl string, uriPaths ...string) (string, error) {
	parseUrl, err := url.Parse(baseUrl)
	if err != nil {
		return "", fmt.Errorf("failed to parse url %s: %w", baseUrl, err)
	}
	var params url.Values
	for _, uriPath := range uriPaths {
		if strings.Contains(uriPath, "?") {
			parts := strings.Split(uriPath, "?")
			parseUrl.Path = path.Join(parseUrl.Path, parts[0])
			params = url.Values{}
			queryParts := strings.Split(parts[1], "=")
			params.Add(queryParts[0], queryParts[1])
			break
		} else {
			parseUrl.Path = path.Join(parseUrl.Path, uriPath)
		}
	}
	if len(params) > 0 {
		parseUrl.RawQuery = params.Encode()
	}
	return parseUrl.String(), nil
}

func getResource(client CloudClient, dest interface{}, uriPaths ...string) error {
	resp, err := client.Get(uriPaths...)
	if err != nil {
		return err
	}
	defer closeResource(resp.Body)
	if resp.StatusCode == http.StatusNotFound {
		return notFoundError
	}
	if resp.StatusCode != http.StatusOK {
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
