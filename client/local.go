package client

import (
	"github.com/Piszmog/cfservices/credentials"
	"github.com/Piszmog/cloudconfigclient/net"
	"github.com/pkg/errors"
	"os"
	"strings"
)

const (
	environmentLocalConfigServerUrls = "CONFIG_SERVER_URLS"
)

func CreateLocalClient() ConfigClient {
	serviceCredentials, _ := GetLocalCredentials()
	configClients := make([]Client, len(serviceCredentials.Credentials))
	for index, cred := range serviceCredentials.Credentials {
		configUri := cred.Uri
		client := net.CreateDefaultHttpClient()
		configClients[index] = Client{configUri: configUri, httpClient: client}
	}
	return ConfigClient{Clients: configClients}
}

func GetLocalCredentials() (*credentials.ServiceCredentials, error) {
	localUrls := os.Getenv(environmentLocalConfigServerUrls)
	if len(localUrls) == 0 {
		return nil, errors.Errorf("No local Config Server URLs provided in environment variable %s", environmentLocalConfigServerUrls)
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
