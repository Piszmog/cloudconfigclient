package configuration

import (
	"encoding/json"
	"github.com/Piszmog/cloudconfigclient/client"
	"github.com/Piszmog/cloudconfigclient/net"
	"github.com/pkg/errors"
)

type Configuration interface {
	GetConfiguration(applicationName string, profiles []string) (*Source, error)
}

type Client struct {
	configClient *client.ConfigClient
}

func CreateconfigurationClient(configClient *client.ConfigClient) *Client {
	return &Client{configClient: configClient}
}

func (client *Client) GetConfiguration(applicationName string, profiles []string) (*Source, error) {
	for _, configClient := range client.configClient.Clients {
		resp, err := configClient.Get(applicationName, net.JoinProfiles(profiles))
		if resp != nil && resp.StatusCode == 404 {
			continue
		}
		if err != nil {
			return nil, errors.Wrapf(err, "failed to retrieve application configurations from %s",
				configClient.GetFullUrl(applicationName, net.JoinProfiles(profiles)))
		}
		if resp.StatusCode != 200 {
			return nil, errors.Errorf("server responded with status code %d from url %s",
				resp.StatusCode,
				configClient.GetFullUrl(applicationName, net.JoinProfiles(profiles)))
		}
		configuration := &Source{}
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(configuration)
		resp.Body.Close()
		if err != nil {
			return nil, errors.Wrapf(err, "failed to decode response from url %s",
				configClient.GetFullUrl(applicationName, net.JoinProfiles(profiles)))
		}
		return configuration, nil
	}
	return nil, errors.Errorf("failed to find configuration for application %s with profiles %s", applicationName, profiles)
}
