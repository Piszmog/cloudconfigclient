package cloudconfigclient

import (
	"errors"
)

const (
	defaultApplicationName    = "default"
	defaultApplicationProfile = "default"
)

var useDefaultLabel = map[string]string{"useDefaultLabel": "true"}

// Resource interface describes how to retrieve files from the Config Server.
type Resource interface {
	// GetFile retrieves the specified file from the provided directory from the Config Server's default branch.
	//
	// The file will be deserialized into the specified interface type.
	GetFile(directory string, file string, interfaceType any) error
	// GetFileFromBranch retrieves the specified file from the provided branch in the provided directory.
	//
	// The file will be deserialized into the specified interface type.
	GetFileFromBranch(branch string, directory string, file string, interfaceType any) error
	// GetFileRaw retrieves the file from the default branch as a byte slice.
	GetFileRaw(directory string, file string) ([]byte, error)
	// GetFileFromBranchRaw retrieves the file from the specified branch as a byte slice.
	GetFileFromBranchRaw(branch string, directory string, file string) ([]byte, error)
}

// GetFile retrieves the specified file from the provided directory from the Config Server's default branch.
//
// The file will be deserialized into the specified interface type.
func (c *Client) GetFile(directory string, file string, interfaceType any) error {
	return c.getFile([]string{defaultApplicationName, defaultApplicationProfile, directory, file}, useDefaultLabel, interfaceType)
}

// GetFileFromBranch retrieves the specified file from the provided branch in the provided directory.
//
// The file will be deserialized into the specified interface type.
func (c *Client) GetFileFromBranch(branch string, directory string, file string, interfaceType any) error {
	return c.getFile([]string{defaultApplicationName, defaultApplicationProfile, branch, directory, file}, nil, interfaceType)
}

func (c *Client) getFile(paths []string, params map[string]string, interfaceType any) error {
	fileFound := false
	for _, client := range c.clients {
		if err := client.GetResource(paths, params, interfaceType); err != nil {
			if errors.Is(err, ErrResourceNotFound) {
				continue
			}
			return err
		}
		fileFound = true
	}
	if !fileFound {
		return errors.New("failed to find file in the Config Server")
	}
	return nil
}

// GetFileRaw retrieves the file from the default branch as a byte slice.
func (c *Client) GetFileRaw(directory string, file string) ([]byte, error) {
	return c.getFileRaw([]string{defaultApplicationName, defaultApplicationProfile, directory, file}, useDefaultLabel)
}

// GetFileFromBranchRaw retrieves the file from the specified branch as a byte slice.
func (c *Client) GetFileFromBranchRaw(branch string, directory string, file string) ([]byte, error) {
	return c.getFileRaw([]string{defaultApplicationName, defaultApplicationProfile, branch, directory, file}, nil)
}

func (c *Client) getFileRaw(paths []string, params map[string]string) (b []byte, err error) {
	fileFound := false
	for _, client := range c.clients {
		b, err = client.GetResourceRaw(paths, params)
		if err != nil {
			if errors.Is(err, ErrResourceNotFound) {
				continue
			}
			return
		}
		fileFound = true
	}
	if !fileFound {
		err = errors.New("failed to find file in the Config Server")
	}
	return
}
