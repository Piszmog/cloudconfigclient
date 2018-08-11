package client

import (
	"encoding/json"
	"github.com/pkg/errors"
)

const (
	defaultApplicationName    = "default"
	defaultApplicationProfile = "default"
)

type Resource interface {
	GetFile(directory string, file string, interfaceType interface{}) error
	GetFileFromBranch(branch string, directory string, file string, interfaceType interface{}) error
}

func (configClient *ConfigClient) GetFile(directory string, file string, interfaceType interface{}) error {
	fileFound := false
	for _, client := range configClient.Clients {
		resp, err := client.Get(defaultApplicationName, defaultApplicationProfile, directory, file+"?useDefaultLabel=true")
		if resp != nil && resp.StatusCode == 404 {
			continue
		}
		if err != nil {
			return errors.Wrapf(err, "failed to retrieve file")
		}
		if resp.StatusCode != 200 {
			return errors.Errorf("server responded with status code %d", resp.StatusCode)
		}
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(interfaceType)
		resp.Body.Close()
		if err != nil {
			return errors.Wrapf(err, "failed to decode response")
		}
		fileFound = true
	}
	if !fileFound {
		return errors.Errorf("failed to find file %s in the Config Server", file)
	}
	return nil
}

func (configClient *ConfigClient) GetFileFromBranch(branch string, directory string, file string, interfaceType interface{}) error {
	fileFound := false
	for _, client := range configClient.Clients {
		resp, err := client.Get(defaultApplicationName, defaultApplicationProfile, branch, directory, file)
		if resp != nil && resp.StatusCode == 404 {
			continue
		}
		if err != nil {
			return errors.Wrapf(err, "failed to retrieve file")
		}
		if resp.StatusCode != 200 {
			return errors.Errorf("server responded with status code %d", resp.StatusCode)
		}
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(interfaceType)
		resp.Body.Close()
		if err != nil {
			return errors.Wrapf(err, "failed to decode response")
		}
		fileFound = true
	}
	if !fileFound {
		return errors.Errorf("failed to find file %s in the Config Server", file)
	}
	return nil
}
