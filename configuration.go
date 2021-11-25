package cloudconfigclient

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"
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
	Data            map[string]interface{}
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
	return PropertySource{}, ErrPropertySourceDoesNotExist
}

// ErrPropertySourceDoesNotExist is the error that is returned when there are no PropertySource that match the specified
// file name.
var ErrPropertySourceDoesNotExist = errors.New("property source does not exist")

// PropertySourceHandler handles the specific PropertySource.
type PropertySourceHandler func(propertySource PropertySource)

// HandlePropertySources handles all PropertySource configurations that are files. This is a convenience method to
// handle boilerplate for-loop code and filtering of non-configuration files.
//
// Config Server may return other configurations (e.g. credhub property sources) that contain no configurations
// (PropertySource.Source is empty).
func (s *Source) HandlePropertySources(handler PropertySourceHandler) {
	for _, propertySource := range s.PropertySources {
		if len(filepath.Ext(propertySource.Name)) > 0 {
			handler(propertySource)
		}
	}
}

func (s *Source) Get(key string, defaultValue string) interface{} {
	for _, propertySource := range s.PropertySources {
		for _key, value := range propertySource.Source {
			fmt.Printf("key=%s, value=%v\n", _key, value)
			if key == _key {
				return value
			}
		}
	}
	return defaultValue
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
	// GetConfiguration retrieves the configurations/property sources of an application based on the name of the application
	// and the profiles of the application.
	GetConfiguration(applicationName string, profiles ...string) (Source, error)
}

// GetConfiguration retrieves the configurations/property sources of an application based on the name of the application
// and the profiles of the application.
func (c *Client) GetConfiguration(applicationName string, profiles []string, label string) (Source, error) {
	var source Source
	paths := []string{applicationName, joinProfiles(profiles)}
	if label != "" {
		paths = append(paths, label)
	}
	for _, client := range c.clients {
		if err := client.GetResource(paths, nil, &source); err != nil {
			if errors.Is(err, ErrResourceNotFound) {
				continue
			}
			return Source{}, err
		}
		source.toMap()
		return source, nil
	}
	return Source{}, fmt.Errorf("failed to find configuration for application %s with profiles %s", applicationName, profiles)
}

func joinProfiles(profiles []string) string {
	return strings.Join(profiles, ",")
}

func (s *Source) toMap() {
	result := map[string]interface{}{}
	for _, propertySource := range s.PropertySources {
		for key, value := range propertySource.Source {
			entries := strings.Split(key, ".")
			result = insertInMapRecursion(entries, value, result)
		}
	}
	s.Data = result
}

func insertInMap(s []string, value interface{}, dest map[string]interface{}) map[string]interface{} {
	keys := s[:len(s)-1]
	last := s[len(s)-1]

	curr := dest
	for _, key := range keys {
		switch curr[key].(type) {
		case nil:
			curr[key] = map[string]interface{}{}
			curr = curr[key].(map[string]interface{})
		case map[string]interface{}:
			curr = curr[key].(map[string]interface{})
		}
	}
	if curr[last] == nil {
		curr[last] = value
	}
	return dest
}

func insertInMapRecursion(s []string, value interface{}, dest map[string]interface{}) map[string]interface{} {
	key := s[0]
	if len(s) > 1 {
		switch dest[key].(type) {
		case nil:
			dest[key] = insertInMap(s[1:], value, map[string]interface{}{})
		case map[string]interface{}:
			dest[key] = insertInMap(s[1:], value, dest[key].(map[string]interface{}))
		}
	} else if len(s) == 1 {
		if dest[key] == nil {
			dest[key] = value
		}
	}
	return dest
}

type BasicAuthTransport struct {
	Username string
	Password string
}

func (bat BasicAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s",
		base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s",
			bat.Username, bat.Password)))))
	return http.DefaultTransport.RoundTrip(req)
}

func (bat *BasicAuthTransport) client() *http.Client {
	return &http.Client{Transport: bat}
}

type ConfigsEnv struct {
	username string
	password string
	url      string
	name     string
	profiles string
	label    string
}

func (c *ConfigsEnv) ReadFromEnv() *ConfigsEnv {
	c.url = os.Getenv("CONFIG_URL")
	c.username = os.Getenv("CONFIG_USERNAME")
	c.password = os.Getenv("CONFIG_PASSWORD")
	c.name = os.Getenv("APPLICATION_NAME")
	c.profiles = os.Getenv("APPLICATION_PROFILES")
	c.label = os.Getenv("CONFIG_LABEL")
	return c
}
func (configsEnv ConfigsEnv) Load() (Source, error) {
	transport := BasicAuthTransport{Username: configsEnv.username, Password: configsEnv.password}
	client := transport.client()
	configConf := Local(client, configsEnv.url)
	configClient, err := New(configConf)

	if err != nil {
		fmt.Println(err)
		return Source{}, err
	}

	// Retrieves the configurations from the Config Server based on the application name, active profiles and label
	source, err := configClient.GetConfiguration(configsEnv.name, strings.Split(configsEnv.profiles, ","), configsEnv.label)
	if err != nil {
		fmt.Println(err)
		return Source{}, err
	}

	return source, nil
}
