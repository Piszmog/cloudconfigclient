package configuration

import (
	"encoding/json"
	"github.com/Piszmog/cloudconfigclient/model"
	"github.com/Piszmog/cloudconfigclient/net"
	"github.com/pkg/errors"
	"net/http"
)

type Configuration interface {
	GetConfiguration(applicationName string, profiles []string) (*model.Configuration, error)
}

type Client struct {
	HttpClient *http.Client
	BaseUrls   []string
}

func CreateClient(urls ...string) *Client {
	return &Client{
		HttpClient: net.CreateDefaultHttpClient(),
		BaseUrls:   urls,
	}
}

func (client *Client) GetConfiguration(applicationName string, profiles []string) (*model.Configuration, error) {
	for _, baseUrl := range client.BaseUrls {
		fullUrl := net.CreateUrl(baseUrl, applicationName, net.JoinProfiles(profiles))
		resp, err := client.HttpClient.Get(fullUrl)
		if resp != nil && resp.StatusCode == 404 {
			continue
		}
		if err != nil {
			return nil, errors.Wrapf(err, "failed to retrieve application configurations from %s", fullUrl)
		}
		if resp.StatusCode != 200 {
			return nil, errors.Errorf("server responded with status code %d from url %s", resp.StatusCode, fullUrl)
		}
		configuration := &model.Configuration{}
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(configuration)
		resp.Body.Close()
		if err != nil {
			return nil, errors.Wrapf(err, "failed to decode response from url %s", fullUrl)
		}
		return configuration, nil
	}
	return nil, errors.Errorf("failed to find configuration for application %s with profiles %s", applicationName, profiles)
}
