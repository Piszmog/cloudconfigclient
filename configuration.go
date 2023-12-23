package cloudconfigclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
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

// Unmarshal converts the Source.PropertySources to the specified type. The type must be a pointer to a struct.
//
// The provided pointer struct must use JSON tags to map to the PropertySource.Source.
//
// This function is not optimized (ugly) and is intended to only be used at startup.
func (s *Source) Unmarshal(v interface{}) error {
	// covert to a map[string]interface{} so we can convert to the target type
	obj, err := toJson(s.PropertySources)
	if err != nil {
		return err
	}
	// convert to bytes
	b, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	// now we can get to our target type
	return json.Unmarshal(b, v)
}

var sliceRegex = regexp.MustCompile(`(.*)\[(\d+)]`)

func toJson(propertySources []PropertySource) (map[string]interface{}, error) {
	// get ready for a wild ride...
	output := map[string]interface{}{}
	// save the root, so we can get back there when we walk the tree
	root := output
	_ = root
	for _, propertySource := range propertySources {
		for k, v := range propertySource.Source {
			keys := strings.Split(k, ".")
			for i, key := range keys {
				// determine if we are detailing with a slice - e.g. foo[0] or bar[0]
				matches := sliceRegex.FindStringSubmatch(key)
				if matches != nil {
					actualKey := matches[1]
					if _, ok := output[actualKey]; !ok {
						output[actualKey] = []interface{}{}
					}
					if len(keys)-1 == i {
						// the value go straight into the slice, we don't have any slice of objects
						output[actualKey] = append(output[actualKey].([]interface{}), v)
						output = root
					} else {
						// ugh... we have a slice of objects
						// convert the index of the path, the index now matters
						index, err := strconv.Atoi(matches[2])
						if err != nil {
							return nil, err
						}
						var obj map[string]interface{}
						slice := output[actualKey].([]interface{})
						// determine if the index we are walking exists yet in the slice we have built up
						if len(slice) > index {
							obj = slice[index].(map[string]interface{})
							if obj == nil {
								obj = map[string]interface{}{}
							}
						} else {
							// the index does not exist, so we need to create it
							for j := len(slice); j <= index; j++ {
								output[actualKey] = append(output[actualKey].([]interface{}), map[string]interface{}{})
							}
							obj = output[actualKey].([]interface{})[index].(map[string]interface{})
						}
						output = obj
					}
				} else if len(keys)-1 == i {
					// the value go straight into the key
					output[key] = v
					output = root
				} else {
					// need to create a nested object
					if _, ok := output[key]; !ok {
						output[key] = map[string]interface{}{}
					}
					output = output[key].(map[string]interface{})
				}
			}
		}
	}
	return output, nil
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
func (c *Client) GetConfiguration(applicationName string, profiles ...string) (Source, error) {
	var source Source
	paths := []string{applicationName, joinProfiles(profiles)}
	for _, client := range c.clients {
		if err := client.GetResource(paths, nil, &source); err != nil {
			if errors.Is(err, ErrResourceNotFound) {
				continue
			}
			return Source{}, err
		}
		return source, nil
	}
	return Source{}, fmt.Errorf("failed to find configuration for application %s with profiles %s", applicationName, profiles)
}

// GetConfiguration retrieves the configurations/property sources of an application based on the name of the application
// and the profiles of the application and the label.
func (c *Client) GetConfigurationWithLabel(label string, applicationName string, profiles ...string) (Source, error) {
	var source Source
	paths := []string{applicationName, joinProfiles(profiles), label}
	for _, client := range c.clients {
		if err := client.GetResource(paths, nil, &source); err != nil {
			if errors.Is(err, ErrResourceNotFound) {
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
