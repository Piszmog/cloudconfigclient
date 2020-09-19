package cloudconfigclient_test

import (
	"github.com/Piszmog/cfservices"
	"github.com/Piszmog/cloudconfigclient"
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
	configClient, err := cloudconfigclient.CreateCloudClient()
	if err != nil {
		t.Errorf("failed to create cloud client with error %v", err)
	}
	if configClient == nil {
		t.Error("failed to create cloud client")
	}
}

func TestCreateCloudClientWhenENVNotSet(t *testing.T) {
	configClient, err := cloudconfigclient.CreateCloudClient()
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
	serviceCredentials, err := cloudconfigclient.GetCloudCredentials(defaultConfigServerName)
	if err != nil {
		t.Errorf("failed to create cloud credentials with error %v", err)
	}
	if serviceCredentials == nil {
		t.Error("failed to create cloud credentials")
	}
}

func TestGetCloudCredentialsWhenENVNotSet(t *testing.T) {
	serviceCredentials, err := cloudconfigclient.GetCloudCredentials(defaultConfigServerName)
	if err == nil {
		t.Error("expected error when env is not set")
	}
	if serviceCredentials != nil {
		t.Error("created cloud credentials when env is not set")
	}
}

func TestCreateOAuth2Client(t *testing.T) {
	credentials := &cfservices.Credentials{
		AccessTokenUri: "tokenUri",
		ClientSecret:   "clientSecret",
		ClientId:       "clientId",
	}
	client, err := cloudconfigclient.CreateOAuth2HTTPClient(credentials)
	if err != nil {
		t.Errorf("failed to create oauth2 client with error %v", err)
	}
	if client == nil {
		t.Error("no oauth2 client returned")
	}
}

func TestCreateOAuth2ClientWhenCredentialsAreNil(t *testing.T) {
	client, err := cloudconfigclient.CreateOAuth2HTTPClient(nil)
	if err == nil {
		t.Error("expected an error when no credentials are passed")
	}
	if client != nil {
		t.Error("able to create an oauth2 client with nil credentials")
	}
}

func TestCreateOauth2Config(t *testing.T) {
	credentials := &cfservices.Credentials{
		AccessTokenUri: "tokenUri",
		ClientSecret:   "clientSecret",
		ClientId:       "clientId",
	}
	config, err := cloudconfigclient.CreateOAuth2HTTPClient(credentials)
	if err != nil {
		t.Errorf("failed to create oauth2 with errpr %v", err)
	}
	if config == nil {
		t.Error("failed to create oauth2 config")
	}
}

func TestCreateOauth2ConfigWhenCredentialsNil(t *testing.T) {
	config, err := cloudconfigclient.CreateOAuth2HTTPClient(nil)
	if err == nil {
		t.Error("expected an error when passing nil credentials when creating oauth2 config")
	}
	if config != nil {
		t.Error("is able to create oauth2 config when credentials are nil")
	}
}
