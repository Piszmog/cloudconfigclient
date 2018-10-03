package cloudconfigclient

import (
	"os"
	"testing"
)

func TestCreateLocalClient(t *testing.T) {
	const localURI = "http://localhost:8080"
	os.Setenv(EnvironmentLocalConfigServerUrls, localURI)
	defer os.Unsetenv(EnvironmentLocalConfigServerUrls)
	configClient, err := CreateLocalClientFromEnv()
	if err != nil {
		t.Errorf("failed to create local client with error %v", err)
	}
	if configClient == nil {
		t.Error("failed to create local client")
	}
}

func TestCreateLocalClientWhenENVNotSet(t *testing.T) {
	configClient, err := CreateLocalClientFromEnv()
	if err == nil {
		t.Errorf("failed to create local client with error %v", err)
	}
	if configClient != nil {
		t.Error("failed to create local client")
	}
}

func TestGetLocalCredentials(t *testing.T) {
	const localURI = "http://localhost:8080"
	os.Setenv(EnvironmentLocalConfigServerUrls, localURI)
	defer os.Unsetenv(EnvironmentLocalConfigServerUrls)
	serviceCredentials, err := GetLocalCredentials()
	if err != nil {
		t.Errorf("failed to get local credentials with error %v", err)
	}
	if serviceCredentials == nil {
		t.Error("failed to create local credentials")
	}
	if serviceCredentials.Credentials[0].Uri != localURI {
		t.Error("local credentials does not have the local url")
	}
}

func TestGetLocalCredentialsWhenEnvNotSet(t *testing.T) {
	serviceCredentials, err := GetLocalCredentials()
	if err == nil {
		t.Errorf("expected an error when creating credentials")
	}
	if serviceCredentials != nil {
		t.Error("created local credentials when uri not set")
	}
}
