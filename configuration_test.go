package cloudconfigclient_test

import (
	"errors"
	"github.com/Piszmog/cloudconfigclient/v2"
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
	mockClient := new(mockCloudClient)
	response := NewMockHttpResponse(200, configurationSource)
	mockClient.On("Get", []string{"appName", "profile"}).Return(response, nil)
	client := NewConfigClient(mockClient)
	_, err := client.GetConfiguration("appName", []string{"profile"})
	assert.NoError(t, err, "failed to retrieve configurations with error")
}

func TestConfigClient_GetConfigurationWhen404(t *testing.T) {
	mockClient := new(mockCloudClient)
	response := NewMockHttpResponse(404, "")
	mockClient.On("Get", []string{"appName", "profile"}).Return(response, nil)
	client := NewConfigClient(mockClient)
	_, err := client.GetConfiguration("appName", []string{"profile"})
	assert.Error(t, err)
}

func TestConfigClient_GetConfigurationWhenError(t *testing.T) {
	mockClient := new(mockCloudClient)
	response := NewMockHttpResponse(500, configurationSource)
	mockClient.On("Get", []string{"appName", "profile"}).Return(response, errors.New("failed"))
	client := NewConfigClient(mockClient)
	_, err := client.GetConfiguration("appName", []string{"profile"})
	assert.Error(t, err)
}

func TestConfigClient_GetConfigurationWhenNoErrorBut500(t *testing.T) {
	mockClient := new(mockCloudClient)
	response := NewMockHttpResponse(500, configurationSource)
	mockClient.On("Get", []string{"appName", "profile"}).Return(response, nil)
	client := NewConfigClient(mockClient)
	_, err := client.GetConfiguration("appName", []string{"profile"})
	assert.Error(t, err)
}

func TestConfigClient_GetConfigurationInvalidResponseBody(t *testing.T) {
	mockClient := new(mockCloudClient)
	response := NewMockHttpResponse(200, "")
	mockClient.On("Get", []string{"appName", "profile"}).Return(response, nil)
	client := NewConfigClient(mockClient)
	_, err := client.GetConfiguration("appName", []string{"profile"})
	assert.Error(t, err)
}

func TestSource_GetPropertySource(t *testing.T) {
	source := cloudconfigclient.Source{
		PropertySources: []cloudconfigclient.PropertySource{
			{Name: "application-foo.yml"},
			{Name: "application-foo.properties"},
			{Name: "test-app-foo.yml"},
		},
	}

	tests := []struct {
		name     string
		fileName string
		found    bool
	}{
		{
			name:     "Property Source Found",
			fileName: "application-foo.yml",
			found:    true,
		},
		{
			name:     "Property Source Not Found - Wrong Extension",
			fileName: "application-foo.json",
			found:    false,
		},
		{
			name:     "Property Source Not Found - Invalid Name",
			fileName: "test.yml",
			found:    false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			propertySource, err := source.GetPropertySource(test.fileName)
			if test.found {
				assert.NoError(t, err)
				assert.Equal(t, test.fileName, propertySource.Name)
			} else {
				assert.ErrorIs(t, err, cloudconfigclient.PropertySourceDoesNotExistErr)
			}
		})
	}
}

func TestSource_HandlePropertySources_NonFileExcluded(t *testing.T) {
	source := cloudconfigclient.Source{
		PropertySources: []cloudconfigclient.PropertySource{
			{Name: "application-foo.yml"},
			{Name: "ssh://foo.bar.com/path/to/repo/path/to/file/application-foo.properties"},
			{Name: "ssh://foo.bar.com/path/to/repo/path/to/file/application-foo.yaml"},
			{Name: "test-app-foo"},
		},
	}
	count := 0
	source.HandlePropertySources(func(propertySource cloudconfigclient.PropertySource) {
		count++
	})
	assert.Equal(t, 3, count)
}
