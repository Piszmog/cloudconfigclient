package client

import (
	"github.com/Piszmog/cfservices"
	"github.com/Piszmog/cloudconfigclient/net"
	"github.com/pkg/errors"
	"os"
	"strings"
)

const (
	// EnvironmentLocalConfigServerUrls is an environment variable for setting base URLs for local Config Servers.
	EnvironmentLocalConfigServerUrls = "CONFIG_SERVER_URLS"
)

// CreateLocalClient creates a ConfigClient for a locally running Config Server.
//
// The ConfigClient's underlying http.Client is configured with timeouts and connection pools.
func CreateLocalClient() (*ConfigClient, error) {
	serviceCredentials, err := GetLocalCredentials()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create a local client")
	}
	configClients := make([]CloudClient, len(serviceCredentials.Credentials))
	for index, cred := range serviceCredentials.Credentials {
		configUri := cred.Uri
		client := net.CreateDefaultHttpClient()
		configClients[index] = Client{configUri: configUri, httpClient: client}
	}
	return &ConfigClient{Clients: configClients}, nil
}

// GetLocalCredentials creates the credentials that are used to configure a ConfigClient to access a local Config Server.
//
// Retrieves the base URLs of Config Servers from the environment variable 'CONFIG_SERVER_URLS' - a comma separated list.
func GetLocalCredentials() (*cfservices.ServiceCredentials, error) {
	localUrls := os.Getenv(EnvironmentLocalConfigServerUrls)
	if len(localUrls) == 0 {
		return nil, errors.Errorf("No local Config Server URLs provided in environment variable %s", EnvironmentLocalConfigServerUrls)
	}
	urls := strings.Split(localUrls, ",")
	creds := make([]cfservices.Credentials, len(urls))
	for index, url := range urls {
		creds[index] = cfservices.Credentials{
			Uri: url,
		}
	}
	return &cfservices.ServiceCredentials{Credentials: creds}, nil
}
