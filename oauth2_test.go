package cloudconfigclient_test

import (
	"github.com/Piszmog/cfservices"
	"github.com/Piszmog/cloudconfigclient"
	"github.com/stretchr/testify/assert"
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

func TestNewCloudClient(t *testing.T) {
	err := os.Setenv(cfservices.VCAPServices, vcapServices)
	assert.NoError(t, err)
	defer os.Unsetenv(cfservices.VCAPServices)
	configClient, err := cloudconfigclient.NewCloudClient()
	assert.NoError(t, err)
	assert.NotNil(t, configClient)
}

func TestNewCloudClientWhenENVNotSet(t *testing.T) {
	configClient, err := cloudconfigclient.NewCloudClient()
	assert.Error(t, err)
	assert.Nil(t, configClient)
}

func TestGetCloudCredentials(t *testing.T) {
	err := os.Setenv(cfservices.VCAPServices, vcapServices)
	assert.NoError(t, err)
	defer os.Unsetenv(cfservices.VCAPServices)
	serviceCredentials, err := cloudconfigclient.GetCloudCredentials(cloudconfigclient.DefaultConfigServerName)
	assert.NoError(t, err)
	assert.NotNil(t, serviceCredentials)
}

func TestGetCloudCredentialsWhenENVNotSet(t *testing.T) {
	serviceCredentials, err := cloudconfigclient.GetCloudCredentials(cloudconfigclient.DefaultConfigServerName)
	assert.Error(t, err)
	assert.Nil(t, serviceCredentials)
}

func TestNewOAuth2Client(t *testing.T) {
	credentials := &cfservices.Credentials{
		AccessTokenUri: "tokenUri",
		ClientSecret:   "clientSecret",
		ClientId:       "clientId",
	}
	client, err := cloudconfigclient.NewOAuth2HTTPClient(credentials)
	assert.NoError(t, err)
	assert.NotNil(t, client)
}

func TestNewOAuth2ClientWhenCredentialsAreNil(t *testing.T) {
	client, err := cloudconfigclient.NewOAuth2HTTPClient(nil)
	assert.Error(t, err)
	assert.Nil(t, client)
}

func TestNewOAuth2Config(t *testing.T) {
	credentials := &cfservices.Credentials{
		AccessTokenUri: "tokenUri",
		ClientSecret:   "clientSecret",
		ClientId:       "clientId",
	}
	config, err := cloudconfigclient.NewOAuth2HTTPClient(credentials)
	assert.NoError(t, err)
	assert.NotNil(t, config)
}

func TestNewOAuth2ConfigWhenCredentialsNil(t *testing.T) {
	config, err := cloudconfigclient.NewOAuth2HTTPClient(nil)
	assert.Error(t, err)
	assert.Nil(t, config)
}
