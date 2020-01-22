package cloudconfigclient

import (
	"errors"
	"fmt"
	"github.com/Piszmog/cloudconfigclient/net"
)

var notFoundError *NotFoundError

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
func (c ConfigClient) GetConfiguration(applicationName string, profiles []string) (Source, error) {
	var source Source
	for _, client := range c.Clients {
		if err := getResource(client, &source, applicationName, net.JoinProfiles(profiles)); err != nil {
			if errors.As(err, &notFoundError) {
				continue
			}
			return Source{}, err
		}
		return source, nil
	}
	return Source{}, fmt.Errorf("failed to find configuration for application %s with profiles %s", applicationName, profiles)
}
