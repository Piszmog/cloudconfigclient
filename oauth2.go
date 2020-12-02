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
	// ConfigServerName the service name of the Config Server in PCF.
	ConfigServerName = "p-config-server"
	// SpringCloudConfigServerName the service name of the Spring Cloud Config Server in PCF.
	SpringCloudConfigServerName = "p.config-server"
)

// NewCloudClient creates a ConfigClient to access Config Servers running in the cloud (specifically Cloud Foundry).
//
// The environment variable 'VCAP_SERVICES' provides a JSON that contains an entry with the key 'p-config-server' (v2.x)
// or 'p.config-server' (v3.x).
// The service 'p-config-server' is search for first. If not found, 'p.config-server' is searched for.
func NewCloudClient() (*ConfigClient, error) {
	service, err := NewCloudClientForService(ConfigServerName)
	if err != nil {
		// if the v2 config server could not be found, let's check if v3 exists
		if errors.Is(err, cfservices.MissingServiceError) {
			service, err = NewCloudClientForService(SpringCloudConfigServerName)
			if err != nil {
				// if still not found, let's return a very specific message to help point developers in right direction
				if errors.Is(err, cfservices.MissingServiceError) {
					return nil, fmt.Errorf("neither %s or %s exist in environment variable 'VCAP_SERVICES'",
						ConfigServerName, SpringCloudConfigServerName)
				}
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return service, err
}

// NewCloudClientForService creates a ConfigClient to access Config Servers running in the cloud (specifically Cloud Foundry).
//
// The environment variables 'VCAP_SERVICES' provides a JSON. The JSON should contain the entry matching the specified name. This
// entry and used to build an OAuth2 client.
func NewCloudClientForService(name string) (*ConfigClient, error) {
	serviceCredentials, err := GetCloudCredentials(name)
	if err != nil {
		return nil, fmt.Errorf("failed to create cloud client: %w", err)
	}
	return NewOAuth2Client(serviceCredentials.Credentials)
}

// NewOAuth2Client creates a ConfigClient to access Config Servers from an array of credentials.
func NewOAuth2Client(credentials []cfservices.Credentials) (*ConfigClient, error) {
	configClients := make([]CloudClient, len(credentials), len(credentials))
	for index, cred := range credentials {
		configUri := cred.Uri
		client, err := NewOAuth2HTTPClient(&cred)
		if err != nil {
			return nil, fmt.Errorf("failed to create oauth2 client for %s: %w", configUri, err)
		}
		configClients[index] = Client{ConfigUri: configUri, HttpClient: client}
	}
	return &ConfigClient{Clients: configClients}, nil
}

// NewOAuth2HTTPClient creates an OAuth2 http.Client from the provided credentials.
func NewOAuth2HTTPClient(cred *cfservices.Credentials) (*http.Client, error) {
	config, err := NewOAuth2Config(cred)
	if err != nil {
		return nil, err
	}
	return config.Client(context.Background()), nil
}

// NewOAuth2Config creates an OAuth2 config from the provided credentials.
func NewOAuth2Config(cred *cfservices.Credentials) (*clientcredentials.Config, error) {
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
		return nil, fmt.Errorf("failed to get credentials for the service %s: %w", name, err)
	}
	return serviceCreds, nil
}
