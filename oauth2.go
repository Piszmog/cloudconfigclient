package cloudconfigclient

import (
	"context"
	"errors"
	"fmt"
	"github.com/Piszmog/cfservices"
	"golang.org/x/oauth2/clientcredentials"
	"net/http"
)

const (
	// DefaultConfigServerName the default name of the config server in PCF
	DefaultConfigServerName = "p-config-server"
)

// CreateCloudClient creates a ConfigClient to access Config Servers running in the cloud (specifically Cloud Foundry).
//
// The environment variables 'VCAP_SERVICES' provides a JSON that contains an entry with the key 'p-config-server'. This
// entry and used to build an OAuth2 client.
func CreateCloudClient() (*ConfigClient, error) {
	return CreateCloudClientForService(DefaultConfigServerName)
}

// CreateCloudClientForService creates a ConfigClient to access Config Servers running in the cloud (specifically Cloud Foundry).
//
// The environment variables 'VCAP_SERVICES' provides a JSON. The JSON should contain the entry matching the specified name. This
// entry and used to build an OAuth2 client.
func CreateCloudClientForService(name string) (*ConfigClient, error) {
	serviceCredentials, err := GetCloudCredentials(name)
	if err != nil {
		return nil, fmt.Errorf("failed to create cloud client: %w", err)
	}
	return CreateOAuth2Client(serviceCredentials.Credentials)
}

// CreateOAuth2Client creates a ConfigClient to access Config Servers from an array of credentials.
func CreateOAuth2Client(credentials []cfservices.Credentials) (*ConfigClient, error) {
	configClients := make([]CloudClient, len(credentials), len(credentials))
	for index, cred := range credentials {
		configUri := cred.Uri
		client, err := CreateOAuth2HTTPClient(&cred)
		if err != nil {
			return nil, fmt.Errorf("failed to create oauth2 client for %s: %w", configUri, err)
		}
		configClients[index] = Client{ConfigUri: configUri, HttpClient: client}
	}
	return &ConfigClient{Clients: configClients}, nil
}

// CreateOAuth2HTTPClient creates an OAuth2 http.Client from the provided credentials.
func CreateOAuth2HTTPClient(cred *cfservices.Credentials) (*http.Client, error) {
	config, err := CreateOAuth2Config(cred)
	if err != nil {
		return nil, err
	}
	return config.Client(context.Background()), nil
}

// CreateOAuth2Config creates an OAuth2 config from the provided credentials.
func CreateOAuth2Config(cred *cfservices.Credentials) (*clientcredentials.Config, error) {
	if cred == nil {
		return nil, errors.New("no credentials provided")
	}
	return &clientcredentials.Config{
		ClientID:     cred.ClientId,
		ClientSecret: cred.ClientSecret,
		TokenURL:     cred.AccessTokenUri,
	}, nil
}

// GetCloudCredentials retrieves the Config Server's credentials so an OAuth2 client can be created.
func GetCloudCredentials(name string) (*cfservices.ServiceCredentials, error) {
	serviceCreds, err := cfservices.GetServiceCredentialsFromEnvironment(name)
	if err != nil {
		return nil, fmt.Errorf("failed to get credentials for the Config Server service %s: %w", name, err)
	}
	return serviceCreds, nil
}
