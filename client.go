package cloudconfigclient

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Piszmog/cfservices"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	// ConfigServerName the service name of the Config Server in PCF.
	ConfigServerName = "p-config-server"
	// EnvironmentLocalConfigServerUrls is an environment variable for setting base URLs for local Config Servers.
	EnvironmentLocalConfigServerUrls = "CONFIG_SERVER_URLS"
	// SpringCloudConfigServerName the service name of the Spring Cloud Config Server in PCF.
	SpringCloudConfigServerName = "p.config-server"
)

// Client contains the clients of the Config Servers.
type Client struct {
	clients []*HTTPClient
}

// New creates a new Client based on the provided options. A Client can be configured to communicate with
// a local Config Server, an OAuth2 Server, and Config Servers in Cloud Foundry.
//
// At least one option must be provided.
func New(options ...Option) (*Client, error) {
	var clients []*HTTPClient
	if len(options) == 0 {
		return nil, errors.New("at least one option must be provided")
	}
	for _, option := range options {
		if err := option(&clients); err != nil {
			return nil, err
		}
	}
	return &Client{clients: clients}, nil
}

// Option creates a slice of httpClients per Config Server instance.
type Option func(*[]*HTTPClient) error

// LocalEnv creates a clients for a locally running Config Servers. The URLs to the Config Servers are acquired from the
// environment variable 'CONFIG_SERVER_URLS'.
func LocalEnv(client *http.Client) Option {
	return func(clients *[]*HTTPClient) error {
		httpClients, err := newLocalClientFromEnv(client)
		if err != nil {
			return err
		}
		*clients = append(*clients, httpClients...)
		return nil
	}
}

func newLocalClientFromEnv(client *http.Client) ([]*HTTPClient, error) {
	localUrls := os.Getenv(EnvironmentLocalConfigServerUrls)
	if len(localUrls) == 0 {
		return nil, fmt.Errorf("no local Config Server URLs provided in environment variable %s", EnvironmentLocalConfigServerUrls)
	}
	return newLocalClient(client, strings.Split(localUrls, ",")), nil
}

// Local creates a clients for a locally running Config Servers.
func Local(client *http.Client, urls ...string) Option {
	return func(clients *[]*HTTPClient) error {
		*clients = append(*clients, newLocalClient(client, urls)...)
		return nil
	}
}

func newLocalClient(client *http.Client, urls []string) []*HTTPClient {
	clients := make([]*HTTPClient, len(urls), len(urls))
	for index, baseUrl := range urls {
		clients[index] = &HTTPClient{BaseURL: baseUrl, Client: client}
	}
	return clients
}

// DefaultCFService creates a clients for each Config Servers the application is bounded to in Cloud Foundry. The
// environment variable 'VCAP_SERVICES' provides a JSON that contains an entry with the key 'p-config-server' (v2.x)
// or 'p.config-server' (v3.x).
//
// The service 'p.config-server' is search for first. If not found, 'p-config-server' is searched for.
func DefaultCFService() Option {
	return func(clients *[]*HTTPClient) error {
		services, err := cfservices.GetServices()
		if err != nil {
			return fmt.Errorf("failed to parse 'VCAP_SERVICES': %w", err)
		}
		httpClients, err := newCloudClientForService(SpringCloudConfigServerName, services)
		if err != nil {
			if errors.Is(err, cfservices.MissingServiceError) {
				httpClients, err = newCloudClientForService(ConfigServerName, services)
				if err != nil {
					if errors.Is(err, cfservices.MissingServiceError) {
						return fmt.Errorf("neither %s or %s exist in environment variable 'VCAP_SERVICES'",
							ConfigServerName, SpringCloudConfigServerName)
					}
					return err
				}
			} else {
				return err
			}
		}
		*clients = append(*clients, httpClients...)
		return nil
	}
}

// CFService creates a clients for each Config Servers the application is bounded to in Cloud Foundry. The environment
// variable 'VCAP_SERVICES' provides a JSON. The JSON should contain the entry matching the specified name. This
// entry and used to build an OAuth Client.
func CFService(service string) Option {
	return func(clients *[]*HTTPClient) error {
		services, err := cfservices.GetServices()
		if err != nil {
			return fmt.Errorf("failed to parse 'VCAP_SERVICES': %w", err)
		}
		httpClients, err := newCloudClientForService(service, services)
		if err != nil {
			return err
		}
		*clients = append(*clients, httpClients...)
		return nil
	}
}

func newCloudClientForService(name string, services map[string][]cfservices.Service) ([]*HTTPClient, error) {
	creds, err := cfservices.GetServiceCredentials(services, name)
	if err != nil {
		return nil, fmt.Errorf("failed to create cloud Client: %w", err)
	}
	clients := make([]*HTTPClient, len(creds.Credentials), len(creds.Credentials))
	for i, cred := range creds.Credentials {
		clients[i] = &HTTPClient{BaseURL: cred.Uri, Client: newOAuth2Client(cred.ClientId, cred.ClientSecret, cred.AccessTokenUri)}
	}
	return clients, nil
}

// OAuth2 creates a Client for a Config Server based on the provided OAuth2.0 information.
func OAuth2(baseURL string, clientId string, secret string, tokenURI string) Option {
	return func(clients *[]*HTTPClient) error {
		*clients = append(*clients, &HTTPClient{BaseURL: baseURL, Client: newOAuth2Client(clientId, secret, tokenURI)})
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
