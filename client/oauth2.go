package client

import (
	"github.com/Piszmog/cfservices"
	"github.com/Piszmog/cloudconfigclient/net"
	"github.com/pkg/errors"
)

const (
	defaultConfigServerName = "p-config-server"
)

func CreateCloudClient() (*ConfigClient, error) {
	return CreateCloudClientForService(defaultConfigServerName)
}

func CreateCloudClientForService(name string) (*ConfigClient, error) {
	serviceCredentials, err := GetCloudCredentials(defaultConfigServerName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create cloud client")
	}
	configClients := make([]Client, len(serviceCredentials.Credentials))
	for index, cred := range serviceCredentials.Credentials {
		configUri := cred.Uri
		client := net.CreateOAuth2Client(cred)
		configClients[index] = Client{configUri: configUri, httpClient: client}
	}
	return &ConfigClient{Clients: configClients}, nil
}

func GetCloudCredentials(name string) (*cfservices.ServiceCredentials, error) {
	vcapServices := cfservices.LoadFromEnvironment()
	serviceCreds, err := cfservices.GetServiceCredentials(name, vcapServices)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get credentials for the Config Server service %s", name)
	}
	return serviceCreds, nil
}
