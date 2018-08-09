package client

import (
	"github.com/Piszmog/cfservices"
	"github.com/Piszmog/cfservices/credentials"
	"github.com/Piszmog/cloudconfigclient/net"
	"github.com/pkg/errors"
)

const (
	DefaultConfigServerName = "p-config-server"
)

func CreateCloudClient() ConfigClient {
	serviceCredentials, _ := GetCloudCredentialsByDefaultName()
	configClients := make([]Client, len(serviceCredentials.Credentials))
	for index, cred := range serviceCredentials.Credentials {
		configUri := cred.Uri
		client := net.CreateOAuth2Client(cred)
		configClients[index] = Client{configUri: configUri, httpClient: client}
	}
	return ConfigClient{Clients: configClients}
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
