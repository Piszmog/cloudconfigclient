package cloudconfigclient

import (
	"errors"
	"fmt"
)

const (
	defaultApplicationName    = "default"
	defaultApplicationProfile = "default"
)

// Resource interface describes how to retrieve files from the Config Server.
type Resource interface {
	GetFile(directory string, file string, interfaceType interface{}) error
	GetFileFromBranch(branch string, directory string, file string, interfaceType interface{}) error
}

// GetFile retrieves the specified file from the provided directory from the Config Server's default branch.
//
// The file will be deserialize into the specified interface type.
func (c ConfigClient) GetFile(directory string, file string, interfaceType interface{}) error {
	fileFound := false
	for _, client := range c.Clients {
		if err := getResource(client, interfaceType, defaultApplicationName, defaultApplicationProfile, directory, file+"?useDefaultLabel=true"); err != nil {
			var notFoundError *NotFoundError
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
func (c *ConfigClient) GetFileFromBranch(branch string, directory string, file string, interfaceType interface{}) error {
	fileFound := false
	for _, client := range c.Clients {
		if err := getResource(client, interfaceType, defaultApplicationName, defaultApplicationProfile, branch, directory, file); err != nil {
			var notFoundError *NotFoundError
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
