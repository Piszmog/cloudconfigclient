package cloudconfigclient_test

import (
	"github.com/Piszmog/cloudconfigclient"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
)

func TestNewLocalClient(t *testing.T) {
	const localURI = "http://localhost:8080"
	err := os.Setenv(cloudconfigclient.EnvironmentLocalConfigServerUrls, localURI)
	assert.NoError(t, err)
	defer os.Unsetenv(cloudconfigclient.EnvironmentLocalConfigServerUrls)
	configClient, err := cloudconfigclient.NewLocalClientFromEnv(&http.Client{})
	assert.NoError(t, err, "failed to create local client with error")
	assert.NotNil(t, configClient)
}

func TestNewLocalClientWhenENVNotSet(t *testing.T) {
	configClient, err := cloudconfigclient.NewLocalClientFromEnv(&http.Client{})
	assert.Error(t, err)
	assert.Nil(t, configClient)
}

func TestGetLocalCredentials(t *testing.T) {
	const localURI = "http://localhost:8080"
	err := os.Setenv(cloudconfigclient.EnvironmentLocalConfigServerUrls, localURI)
	assert.NoError(t, err)
	defer os.Unsetenv(cloudconfigclient.EnvironmentLocalConfigServerUrls)
	serviceCredentials, err := cloudconfigclient.GetLocalCredentials()
	assert.NoError(t, err, "failed to get local credentials with error")
	assert.NotNil(t, serviceCredentials)
	assert.Equal(t, localURI, serviceCredentials.Credentials[0].Uri)
}

func TestGetLocalCredentialsWhenEnvNotSet(t *testing.T) {
	serviceCredentials, err := cloudconfigclient.GetLocalCredentials()
	assert.Error(t, err)
	assert.Nil(t, serviceCredentials)
}
