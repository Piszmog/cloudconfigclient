package cloudconfigclient

import (
	"encoding/json"
	"github.com/Piszmog/cloudconfigclient/net"
	"github.com/pkg/errors"
)

// Source is the application's source configurations. It con contain zero to n number of property sources.
type Source struct {
	Name            string           `json:"name"`
	Profiles        []string         `json:"profiles"`
	Label           string           `json:"label"`
	Version         string           `json:"version"`
	State           string           `json:"state"`
	PropertySources []PropertySource `json:"propertySources"`
}

// PropertySource is the property source for the application.
//
// A property source is either a YAML or a PROPERTIES file located in the repository that a Config Server is pointed at.
type PropertySource struct {
	Name   string                 `json:"name"`
	Source map[string]interface{} `json:"source"`
}

// Configuration interface for retrieving an application's configuration files from the Config Server.
type Configuration interface {
	GetConfiguration(applicationName string, profiles []string) (*Source, error)
}

// GetConfiguration retrieves the configurations/property sources of an application based on the name of the application
// and the profiles of the application.
func (configClient ConfigClient) GetConfiguration(applicationName string, profiles []string) (*Source, error) {
	for _, client := range configClient.Clients {
		resp, err := client.Get(applicationName, net.JoinProfiles(profiles))
		if resp != nil && resp.StatusCode == 404 {
			continue
		}
		if err != nil {
			return nil, errors.Wrap(err, "failed to retrieve application configurations")
		}
		if resp.StatusCode != 200 {
			return nil, errors.Errorf("server responded with status code %d", resp.StatusCode)
		}
		configuration := &Source{}
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(configuration)
		resp.Body.Close()
		if err != nil {
			return nil, errors.Wrap(err, "failed to decode response from url")
		}
		return configuration, nil
	}
	return nil, errors.Errorf("failed to find configuration for application %s with profiles %s", applicationName, profiles)
}
