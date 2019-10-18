package cloudconfigclient

import (
	"errors"
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
	configClient := createMockConfigClient(200, configurationSource, nil)
	_, err := configClient.GetConfiguration("appName", []string{"profile"})
	if err != nil {
		t.Errorf("failed to retrieve configurations with error %v", err)
	}
}

func TestConfigClient_GetConfigurationWhen404(t *testing.T) {
	configClient := createMockConfigClient(404, "", nil)
	_, err := configClient.GetConfiguration("appName", []string{"profile"})
	if err == nil {
		t.Error("expected an error to occur")
	}
}

func TestConfigClient_GetConfigurationWhenError(t *testing.T) {
	configClient := createMockConfigClient(500, "", errors.New("failed"))
	_, err := configClient.GetConfiguration("appName", []string{"profile"})
	if err == nil {
		t.Error("expected an error to occur")
	}
}

func TestConfigClient_GetConfigurationWhenNoErrorBut500(t *testing.T) {
	configClient := createMockConfigClient(500, "", nil)
	_, err := configClient.GetConfiguration("appName", []string{"profile"})
	if err == nil {
		t.Error("expected an error to occur")
	}
}

func TestConfigClient_GetConfigurationInvalidResponseBody(t *testing.T) {
	configClient := createMockConfigClient(200, "", nil)
	_, err := configClient.GetConfiguration("appName", []string{"profile"})
	if err == nil {
		t.Error("expected an error to occur")
	}
}
