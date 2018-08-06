package main

import (
	"fmt"
	"github.com/Piszmog/cfservices"
	"github.com/Piszmog/cfservices/credentials"
	"github.com/Piszmog/cloudconfigclient/configuration"
	"github.com/Piszmog/cloudconfigclient/resource"
	"github.com/pkg/errors"
	"os"
	"strings"
)

const (
	DefaultConfigServerName          = "p-config-server"
	EnvironmentLocalConfigServerUrls = "CONFIG_SERVER_URLS"
)

// just for testing
type File struct {
	Example Example `json:"example"`
}

type Example struct {
	Field3 string `json:"field3"`
}

// just for testing -- remove after library built out
func main() {
	serviceCreds, err := GetLocalCredentials()
	if err != nil {
		panic(err)
	}
	var urls []string
	for _, cred := range serviceCreds.Credentials {
		urls = append(urls, cred.Uri)
	}
	file := &File{}
	resourceClient := resource.CreateClient(urls...)
	err = resourceClient.GetFileFromBranch("develop", "temp", "temp1.json", file)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", file)

	configClient := configuration.CreateClient(urls...)
	configurations, err := configClient.GetConfiguration("exampleapp", []string{"dev"})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", configurations)
}

func GetLocalCredentials() (*credentials.ServiceCredentials, error) {
	localUrls := os.Getenv(EnvironmentLocalConfigServerUrls)
	if len(localUrls) == 0 {
		return nil, errors.Errorf("No local Config Server URLs provided in environment variable %s", EnvironmentLocalConfigServerUrls)
	}
	urls := strings.Split(localUrls, ",")
	var creds []credentials.Credentials
	for _, url := range urls {
		creds = append(creds, credentials.Credentials{
			Uri: url,
		})
	}
	return &credentials.ServiceCredentials{Credentials: creds}, nil
}

func GetCloudCredentialsByDefaultName() (*credentials.ServiceCredentials, error) {
	return GetCloudCredentials(DefaultConfigServerName)
}

func GetCloudCredentials(name string) (*credentials.ServiceCredentials, error) {
	vcapServices := cfservices.LoadFromEnvironment()
	serviceCreds, err := cfservices.GetServiceCredentials(name, vcapServices)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get credentials for the Config Server")
	}
	return serviceCreds, nil
}
