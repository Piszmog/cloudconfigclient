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
	mockClient.On("Get", []string{"appName", "profile"}).Return(response, nil)
	client := NewConfigClient(mockClient)
	_, err := client.GetConfiguration("appName", []string{"profile"})
	assert.NoError(t, err, "failed to retrieve configurations with error")
}

func TestConfigClient_GetConfigurationWhen404(t *testing.T) {
	mockClient := new(MockCloudClient)
	response := NewMockHttpResponse(404, "")
	mockClient.On("Get", []string{"appName", "profile"}).Return(response, nil)
	client := NewConfigClient(mockClient)
	_, err := client.GetConfiguration("appName", []string{"profile"})
	assert.Error(t, err)
}

func TestConfigClient_GetConfigurationWhenError(t *testing.T) {
	mockClient := new(MockCloudClient)
	response := NewMockHttpResponse(500, configurationSource)
	mockClient.On("Get", []string{"appName", "profile"}).Return(response, errors.New("failed"))
	client := NewConfigClient(mockClient)
	_, err := client.GetConfiguration("appName", []string{"profile"})
	assert.Error(t, err)
}

func TestConfigClient_GetConfigurationWhenNoErrorBut500(t *testing.T) {
	mockClient := new(MockCloudClient)
	response := NewMockHttpResponse(500, configurationSource)
	mockClient.On("Get", []string{"appName", "profile"}).Return(response, nil)
	client := NewConfigClient(mockClient)
	_, err := client.GetConfiguration("appName", []string{"profile"})
	assert.Error(t, err)
}

func TestConfigClient_GetConfigurationInvalidResponseBody(t *testing.T) {
	mockClient := new(MockCloudClient)
	response := NewMockHttpResponse(200, "")
	mockClient.On("Get", []string{"appName", "profile"}).Return(response, nil)
	client := NewConfigClient(mockClient)
	_, err := client.GetConfiguration("appName", []string{"profile"})
	assert.Error(t, err)
}
