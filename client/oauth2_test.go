package client

import (
	"github.com/Piszmog/cfservices"
	"os"
	"testing"
)

const (
	vcapServices = `{
  "p-config-server": [
    {
      "name": "config-server",
      "instance_name": "config-server",
      "binding_name": null,
      "credentials": {
        "uri": "https://config-uri.com",
        "client_secret": "clientSecret",
        "client_id": "config-client-id",
        "access_token_uri": "https://tokenuri.com"
      },
      "syslog_drain_url": null,
      "volume_mounts": [],
      "label": "p-config-server",
      "provider": null,
      "plan": "testPlan",
      "tags": [
        "testTag"
      ]
    }
  ]
}`
)

func TestCreateCloudClient(t *testing.T) {
	os.Setenv(cfservices.VCAPServices, vcapServices)
	defer os.Unsetenv(cfservices.VCAPServices)
	configClient, err := CreateCloudClient()
	if err != nil {
		t.Errorf("failed to create cloud client with error %v", err)
	}
	if configClient == nil {
		t.Error("failed to create cloud client")
	}
}

func TestCreateCloudClientWhenENVNotSet(t *testing.T) {
	configClient, err := CreateCloudClient()
	if err == nil {
		t.Error("expected error when env is not set")
	}
	if configClient != nil {
		t.Error("created cloud client when env is not set")
	}
}

func TestGetCloudCredentials(t *testing.T) {
	os.Setenv(cfservices.VCAPServices, vcapServices)
	defer os.Unsetenv(cfservices.VCAPServices)
	serviceCredentials, err := GetCloudCredentials(defaultConfigServerName)
	if err != nil {
		t.Errorf("failed to create cloud credentials with error %v", err)
	}
	if serviceCredentials == nil {
		t.Error("failed to create cloud credentials")
	}
}

func TestGetCloudCredentialsWhenENVNotSet(t *testing.T) {
	serviceCredentials, err := GetCloudCredentials(defaultConfigServerName)
	if err == nil {
		t.Error("expected error when env is not set")
	}
	if serviceCredentials != nil {
		t.Error("created cloud credentials when env is not set")
	}
}
