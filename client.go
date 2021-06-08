package cloudconfigclient

import (
	"context"
	"errors"
	"fmt"
	"github.com/Piszmog/cfservices"
	"golang.org/x/oauth2/clientcredentials"
	"net/http"
	"os"
	"strings"
)

const (
	// ConfigServerName the service name of the Config Server in PCF.
	ConfigServerName = "p-config-server"
	// EnvironmentLocalConfigServerUrls is an environment variable for setting base URLs for local Config Servers.
	EnvironmentLocalConfigServerUrls = "CONFIG_SERVER_URLS"
	// SpringCloudConfigServerName the service name of the Spring Cloud Config Server in PCF.
	SpringCloudConfigServerName = "p.config-server"
)

// ConfigClient contains the clients of the Config Servers.
type ConfigClient struct {
	// Clients are all the config server clients.
	Clients []*httpClient
}

func New(options ...Option) (*ConfigClient, error) {
	var clients []*httpClient
	if len(options) == 0 {
		return nil, errors.New("at least one option must be provided")
	}
	for _, option := range options {
		if err := option(clients); err != nil {
			return nil, err
		}
	}
	return &ConfigClient{Clients: clients}, nil
}

type Option func([]*httpClient) error

func LocalEnv(client *http.Client) Option {
	return func(clients []*httpClient) error {
		httpClients, err := newLocalClientFromEnv(client)
		if err != nil {
			return err
		}
		clients = append(clients, httpClients...)
		return nil
	}
}

func newLocalClientFromEnv(client *http.Client) ([]*httpClient, error) {
	localUrls := os.Getenv(EnvironmentLocalConfigServerUrls)
	if len(localUrls) == 0 {
		return nil, fmt.Errorf("no local Config Server URLs provided in environment variable %s", EnvironmentLocalConfigServerUrls)
	}
	return newLocalClient(client, strings.Split(localUrls, ",")), nil
}

func Local(client *http.Client, urls []string) Option {
	return func(clients []*httpClient) error {
		clients = append(clients, newLocalClient(client, urls)...)
		return nil
	}
}

func newLocalClient(client *http.Client, urls []string) []*httpClient {
	clients := make([]*httpClient, len(urls), len(urls))
	for index, baseUrl := range urls {
		clients[index] = &httpClient{baseURL: baseUrl, client: client}
	}
	return clients
}

func DefaultCFService() Option {
	return func(clients []*httpClient) error {
		httpClients, err := newCloudClientForService(SpringCloudConfigServerName)
		if err != nil {
			if errors.Is(err, cfservices.MissingServiceError) {
				httpClients, err = newCloudClientForService(SpringCloudConfigServerName)
				if err != nil {
					if errors.Is(err, cfservices.MissingServiceError) {
						return fmt.Errorf("neither %s or %s exist in environment variable 'VCAP_SERVICES'",
							ConfigServerName, SpringCloudConfigServerName)
					} else {
						return err
					}
				}
				clients = append(clients, httpClients...)
			} else {
				return err
			}
		}
		clients = append(clients, httpClients...)
		return nil
	}
}

func CFService(service string) Option {
	return func(clients []*httpClient) error {
		httpClients, err := newCloudClientForService(service)
		if err != nil {
			return err
		}
		clients = append(clients, httpClients...)
		return nil
	}
}

func newCloudClientForService(name string) ([]*httpClient, error) {
	creds, err := getCloudCredentials(name)
	if err != nil {
		return nil, fmt.Errorf("failed to create cloud client: %w", err)
	}
	clients := make([]*httpClient, len(creds.Credentials), len(creds.Credentials))
	for i, cred := range creds.Credentials {
		clients[i] = &httpClient{baseURL: cred.Uri, client: newOAuth2Client(cred.ClientId, cred.ClientSecret, cred.AccessTokenUri)}
	}
	return clients, nil
}

func getCloudCredentials(name string) (*cfservices.ServiceCredentials, error) {
	serviceCreds, err := cfservices.GetServiceCredentialsFromEnvironment(name)
	if err != nil {
		return nil, fmt.Errorf("failed to get credentials for the service %s: %w", name, err)
	}
	return serviceCreds, nil
}

func OAuth2(baseURL string, clientId string, secret string, tokenURI string) Option {
	return func(clients []*httpClient) error {
		clients = append(clients, &httpClient{baseURL: baseURL, client: newOAuth2Client(clientId, secret, tokenURI)})
		return nil
	}
}

func newOAuth2Client(clientId string, secret string, tokenURI string) *http.Client {
	config := newOAuth2Config(clientId, secret, tokenURI)
	return config.Client(context.Background())
}

func newOAuth2Config(clientId string, secret string, tokenURI string) *clientcredentials.Config {
	return &clientcredentials.Config{ClientID: clientId, ClientSecret: secret, TokenURL: tokenURI}
}
