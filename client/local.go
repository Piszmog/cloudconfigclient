package client

import (
	"github.com/Piszmog/cfservices"
	"github.com/Piszmog/cloudconfigclient/net"
	"github.com/pkg/errors"
	"os"
	"strings"
)

const (
	EnvironmentLocalConfigServerUrls = "CONFIG_SERVER_URLS"
)

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

func GetLocalCredentials() (*cfservices.ServiceCredentials, error) {
	localUrls := os.Getenv(EnvironmentLocalConfigServerUrls)
	if len(localUrls) == 0 {
		return nil, errors.Errorf("No local Config Server URLs provided in environment variable %s", EnvironmentLocalConfigServerUrls)
	}
	urls := strings.Split(localUrls, ",")
	var creds []cfservices.Credentials
	for _, url := range urls {
		creds = append(creds, cfservices.Credentials{
			Uri: url,
		})
	}
	return &cfservices.ServiceCredentials{Credentials: creds}, nil
}
