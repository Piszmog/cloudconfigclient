package client

import (
	"encoding/json"
	"github.com/Piszmog/cloudconfigclient/net"
	"github.com/pkg/errors"
)

type Source struct {
	Name            string           `json:"name"`
	Profiles        []string         `json:"profiles"`
	Label           string           `json:"label"`
	Version         string           `json:"version"`
	State           string           `json:"state"`
	PropertySources []PropertySource `json:"propertySources"`
}

type PropertySource struct {
	Name   string                 `json:"name"`
	Source map[string]interface{} `json:"source"`
}

type Configuration interface {
	GetConfiguration(applicationName string, profiles []string) (*Source, error)
}

func (configClient *ConfigClient) GetConfiguration(applicationName string, profiles []string) (*Source, error) {
	for _, client := range configClient.Clients {
		resp, err := client.Get(applicationName, net.JoinProfiles(profiles))
		if resp != nil && resp.StatusCode == 404 {
			continue
		}
		if err != nil {
			return nil, errors.Wrapf(err, "failed to retrieve application configurations")
		}
		if resp.StatusCode != 200 {
			return nil, errors.Errorf("server responded with status code %d", resp.StatusCode)
		}
		configuration := &Source{}
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(configuration)
		resp.Body.Close()
		if err != nil {
			return nil, errors.Wrapf(err, "failed to decode response from url")
		}
		return configuration, nil
	}
	return nil, errors.Errorf("failed to find configuration for application %s with profiles %s", applicationName, profiles)
}
