package cloudconfigclient

import (
	"errors"
	"fmt"
)

const (
	defaultApplicationName    = "default"
	defaultApplicationProfile = "default"
)

var useDefaultLabel = map[string]string{"useDefaultLabel": "true"}

// Resource interface describes how to retrieve files from the Config Server.
type Resource interface {
	GetFile(directory string, file string, interfaceType interface{}) error
	GetFileFromBranch(branch string, directory string, file string, interfaceType interface{}) error
}

// GetFile retrieves the specified file from the provided directory from the Config Server's default branch.
//
// The file will be deserialize into the specified interface type.
func (c *Client) GetFile(directory string, file string, interfaceType interface{}) error {
	fileFound := false
	paths := []string{defaultApplicationName, defaultApplicationProfile, directory, file}
	for _, client := range c.clients {
		if err := client.GetResource(paths, useDefaultLabel, interfaceType); err != nil {
			if errors.As(err, &notFoundError) {
				continue
			}
			return err
		}
		fileFound = true
	}
	if !fileFound {
		return fmt.Errorf("failed to find file %s in the Config Server", file)
	}
	return nil
}

// GetFileFromBranch retrieves the specified file from the provided branch in the provided directory.
//
// The file will be deserialize into the specified interface type.
func (c *Client) GetFileFromBranch(branch string, directory string, file string, interfaceType interface{}) error {
	fileFound := false
	paths := []string{defaultApplicationName, defaultApplicationProfile, branch, directory, file}
	for _, client := range c.clients {
		if err := client.GetResource(paths, nil, interfaceType); err != nil {
			if errors.As(err, &notFoundError) {
				continue
			}
			return err
		}
		fileFound = true
	}
	if !fileFound {
		return fmt.Errorf("failed to find file %s in the Config Server", file)
	}
	return nil
}
