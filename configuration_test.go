package cloudconfigclient_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	configurationSource = `{
  "name": "testConfig",
  "profiles": [
    "profile"
  ],
  "propertySources": [
    {
      "name": "test",
      "source": {
        "field1": "value1",
        "field2": 1
      }
    }
  ]
}`
)

func TestConfigClient_GetConfiguration(t *testing.T) {
	mockClient := new(MockCloudClient)
	response := NewMockHttpResponse(200, configurationSource)
	mockClient.On("Get", "profile").Return(response, nil)
	configClient := createMockConfigClient(200, configurationSource, nil)
	_, err := configClient.GetConfiguration("appName", []string{"profile"})
	assert.NoError(t, err, "failed to retrieve configurations with error")
}

func TestConfigClient_GetConfigurationWhen404(t *testing.T) {
	configClient := createMockConfigClient(404, "", nil)
	_, err := configClient.GetConfiguration("appName", []string{"profile"})
	assert.Error(t, err)
}

func TestConfigClient_GetConfigurationWhenError(t *testing.T) {
	configClient := createMockConfigClient(500, "", errors.New("failed"))
	_, err := configClient.GetConfiguration("appName", []string{"profile"})
	assert.Error(t, err)
}

func TestConfigClient_GetConfigurationWhenNoErrorBut500(t *testing.T) {
	configClient := createMockConfigClient(500, "", nil)
	_, err := configClient.GetConfiguration("appName", []string{"profile"})
	assert.Error(t, err)
}

func TestConfigClient_GetConfigurationInvalidResponseBody(t *testing.T) {
	configClient := createMockConfigClient(200, "", nil)
	_, err := configClient.GetConfiguration("appName", []string{"profile"})
	assert.Error(t, err)
}
