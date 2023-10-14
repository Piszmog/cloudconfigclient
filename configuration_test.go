package cloudconfigclient_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/Piszmog/cloudconfigclient/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestClient_GetConfiguration(t *testing.T) {
	tests := []struct {
		name        string
		application string
		profiles    []string
		checker     func(*testing.T, *http.Request)
		response    *http.Response
		expected    cloudconfigclient.Source
		err         error
	}{
		{
			name:        "Get Config",
			application: "appName",
			profiles:    []string{"profile"},
			checker: func(t *testing.T, request *http.Request) {
				require.Equal(t, "http://localhost:8888/appName/profile", request.URL.String())
			},
			response: NewMockHttpResponse(http.StatusOK, configurationSource),
			expected: cloudconfigclient.Source{
				Name:            "testConfig",
				Profiles:        []string{"profile"},
				PropertySources: []cloudconfigclient.PropertySource{{Name: "test", Source: map[string]interface{}{"field1": "value1", "field2": float64(1)}}},
			},
		},
		{
			name:        "Multiple Profiles",
			application: "appName",
			profiles:    []string{"profile1", "profile2", "profile3"},
			checker: func(t *testing.T, request *http.Request) {
				require.Equal(t, "http://localhost:8888/appName/profile1,profile2,profile3", request.URL.String())
			},
			response: NewMockHttpResponse(http.StatusOK, configurationSource),
			expected: cloudconfigclient.Source{
				Name:            "testConfig",
				Profiles:        []string{"profile"},
				PropertySources: []cloudconfigclient.PropertySource{{Name: "test", Source: map[string]interface{}{"field1": "value1", "field2": float64(1)}}},
			},
		},
		{
			name:        "Not Found",
			application: "appName",
			profiles:    []string{"profile"},
			response:    NewMockHttpResponse(http.StatusNotFound, ""),
			err:         errors.New("failed to find configuration for application appName with profiles [profile]"),
		},
		{
			name:        "Server Error",
			application: "appName",
			profiles:    []string{"profile"},
			response:    NewMockHttpResponse(http.StatusInternalServerError, ""),
			err:         errors.New("server responded with status code '500' and body ''"),
		},
		{
			name:        "No Response Body",
			application: "appName",
			profiles:    []string{"profile"},
			response:    NewMockHttpResponse(http.StatusOK, ""),
			err:         errors.New("failed to decode response from url: EOF"),
		},
		{
			name:        "HTTP Error",
			application: "appName",
			profiles:    []string{"profile"},
			err:         errors.New("failed to retrieve from http://localhost:8888/appName/profile: Get \"http://localhost:8888/appName/profile\": http: RoundTripper implementation (cloudconfigclient_test.RoundTripFunc) returned a nil *Response with a nil error"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			httpClient := NewMockHttpClient(func(req *http.Request) *http.Response {
				if test.checker != nil {
					test.checker(t, req)
				}
				return test.response
			})
			client, err := cloudconfigclient.New(cloudconfigclient.Local(httpClient, &cloudconfigclient.BasicAuthCredential{}, "http://localhost:8888"))
			require.NoError(t, err)
			configuration, err := client.GetConfiguration(test.application, test.profiles...)
			if err != nil {
				require.Error(t, err)
				require.Equal(t, test.err.Error(), err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expected, configuration)
			}
		})
	}
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
				assert.ErrorIs(t, err, cloudconfigclient.ErrPropertySourceDoesNotExist)
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

func TestSource_Unmarshal(t *testing.T) {
	tests := []struct {
		name     string
		source   cloudconfigclient.Source
		expected testStruct
		hasError bool
	}{
		{
			name: "Full Struct",
			source: cloudconfigclient.Source{
				PropertySources: []cloudconfigclient.PropertySource{
					{
						Name: "application-foo.yml",
						Source: map[string]interface{}{
							"stringVal":                "value1",
							"intVal":                   1,
							"sliceString[0]":           "value2",
							"sliceInt[0]":              2,
							"sliceStruct[0].stringVal": "value5",
							"sliceStruct[0].intVal":    4,
							"sliceStruct[1].intVal":    5,
							"sliceStruct[1].stringVal": "value6",
						},
					},
				},
			},
			expected: testStruct{
				StringVal: "value1",
				IntVal:    1,
				SliceString: []string{
					"value2",
				},
				SliceInt: []int{
					2,
				},
				SliceStruct: []nestedStruct{
					{
						StringVal: "value5",
						IntVal:    4,
					},
					{
						StringVal: "value6",
						IntVal:    5,
					},
				},
			},
		},
		{
			name: "Split Across Files",
			source: cloudconfigclient.Source{
				PropertySources: []cloudconfigclient.PropertySource{
					{
						Name: "application-foo.yml",
						Source: map[string]interface{}{
							"stringVal":      "value1",
							"intVal":         1,
							"sliceString[0]": "value2",
							"sliceInt[0]":    2,
						},
					},
					{
						Name: "application-foo-dev.yml",
						Source: map[string]interface{}{
							"sliceStruct[0].stringVal": "value5",
							"sliceStruct[0].intVal":    4,
							"sliceStruct[1].intVal":    5,
							"sliceStruct[1].stringVal": "value6",
						},
					},
				},
			},
			expected: testStruct{
				StringVal: "value1",
				IntVal:    1,
				SliceString: []string{
					"value2",
				},
				SliceInt: []int{
					2,
				},
				SliceStruct: []nestedStruct{
					{
						StringVal: "value5",
						IntVal:    4,
					},
					{
						StringVal: "value6",
						IntVal:    5,
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var actual testStruct
			err := test.source.Unmarshal(&actual)
			if test.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, actual)
			}
		})
	}
}

type testStruct struct {
	StringVal   string         `json:"stringVal"`
	IntVal      int            `json:"intVal"`
	SliceString []string       `json:"sliceString"`
	SliceInt    []int          `json:"sliceInt"`
	SliceStruct []nestedStruct `json:"sliceStruct"`
}

type nestedStruct struct {
	StringVal string `json:"stringVal"`
	IntVal    int    `json:"intVal"`
}
