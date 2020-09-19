package cloudconfigclient_test

import (
	"fmt"
	"github.com/Piszmog/cloudconfigclient"
	"net/http"
	"os"
	"testing"
)

func TestCreateLocalClient(t *testing.T) {
	const localURI = "http://localhost:8080"
	if err := os.Setenv(cloudconfigclient.EnvironmentLocalConfigServerUrls, localURI); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Unsetenv(cloudconfigclient.EnvironmentLocalConfigServerUrls); err != nil {
			fmt.Println(err)
		}
	}()
	configClient, err := cloudconfigclient.NewLocalClientFromEnv(&http.Client{})
	if err != nil {
		t.Errorf("failed to create local client with error %v", err)
	}
	if configClient == nil {
		t.Error("failed to create local client")
	}
}

func TestCreateLocalClientWhenENVNotSet(t *testing.T) {
	configClient, err := cloudconfigclient.NewLocalClientFromEnv(&http.Client{})
	if err == nil {
		t.Errorf("failed to create local client with error %v", err)
	}
	if configClient != nil {
		t.Error("failed to create local client")
	}
}

func TestGetLocalCredentials(t *testing.T) {
	const localURI = "http://localhost:8080"
	if err := os.Setenv(cloudconfigclient.EnvironmentLocalConfigServerUrls, localURI); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Unsetenv(cloudconfigclient.EnvironmentLocalConfigServerUrls); err != nil {
			fmt.Println(err)
		}
	}()
	serviceCredentials, err := cloudconfigclient.GetLocalCredentials()
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
	serviceCredentials, err := cloudconfigclient.GetLocalCredentials()
	if err == nil {
		t.Errorf("expected an error when creating credentials")
	}
	if serviceCredentials != nil {
		t.Error("created local credentials when uri not set")
	}
}
