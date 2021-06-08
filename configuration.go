package cloudconfigclient

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
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

// GetPropertySource retrieves the PropertySource that has the specifies fileName. The fileName is the name of the file
// with extension - e.g. application-foo.yml.
//
// Usually the Config Server will return the PropertySource.Name as an URL of sorts
// (e.g. ssh://base-url.com/path/to/repository/path/to/file.[yml/properties]). So in order to find the specific file
// with the desired configurations, the ending of the name needs to be matched against.
func (s *Source) GetPropertySource(fileName string) (PropertySource, error) {
	for _, propertySource := range s.PropertySources {
		if strings.HasSuffix(propertySource.Name, fileName) {
			return propertySource, nil
		}
	}
	return PropertySource{}, PropertySourceDoesNotExistErr
}

// PropertySourceDoesNotExistErr is the error that is returned when there are no PropertySource that match the specified
// file name.
var PropertySourceDoesNotExistErr = errors.New("property source does not exist")

// PropertySourceHandler handles the specific PropertySource.
type PropertySourceHandler func(propertySource PropertySource)

// HandlePropertySources handles all PropertySource configurations that are files. This is a convenience method to
// handle boilerplate for-loop code and filtering of non-configuration files.
//
// Config Server may return other configurations (e.g. credhub proeprty sources) that contain no configurations
// (PropertySource.Source is empty).
func (s *Source) HandlePropertySources(handler PropertySourceHandler) {
	for _, propertySource := range s.PropertySources {
		if len(filepath.Ext(propertySource.Name)) > 0 {
			handler(propertySource)
		}
	}
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
	GetConfiguration(applicationName string, profiles ...string) (Source, error)
}

// GetConfiguration retrieves the configurations/property sources of an application based on the name of the application
// and the profiles of the application.
func (c *Client) GetConfiguration(applicationName string, profiles ...string) (Source, error) {
	var source Source
	paths := []string{applicationName, joinProfiles(profiles)}
	for _, client := range c.clients {
		if err := client.getResource(paths, nil, &source); err != nil {
			if errors.As(err, &notFoundError) {
				continue
			}
			return Source{}, err
		}
		return source, nil
	}
	return Source{}, fmt.Errorf("failed to find configuration for application %s with profiles %s", applicationName, profiles)
}

func joinProfiles(profiles []string) string {
	return strings.Join(profiles, ",")
}
