package cloudconfigclient_test

import (
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/duvalhub/cloudconfigclient"
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
			client, err := cloudconfigclient.New(cloudconfigclient.Local(httpClient, "http://localhost:8888"))
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

func TestInsertInMap(t *testing.T) {

	type args struct {
		s     []string
		value string
		dest  map[string]interface{}
	}

	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "allo",
			args: args{
				s:     []string{"a"},
				value: "allo",
				dest: map[string]interface{}{
					"a": "toto",
				},
			},
			want: map[string]interface{}{
				"a": "toto",
			},
		},
		{
			name: "allo",
			args: args{
				s:     []string{"a"},
				value: "allo",
				dest: map[string]interface{}{
					"asd": "qwe",
				},
			},
			want: map[string]interface{}{
				"a":   "allo",
				"asd": "qwe",
			},
		},
		{
			name: "b.a",
			args: args{
				s:     []string{"b", "a"},
				value: "allo",
				dest:  map[string]interface{}{},
			},
			want: map[string]interface{}{
				"b": map[string]interface{}{
					"a": "allo",
				},
			},
		},
		{
			name: "b.a already set",
			args: args{
				s:     []string{"b", "a"},
				value: "allo",
				dest: map[string]interface{}{
					"b": map[string]interface{}{
						"a": "bye",
					},
				},
			},
			want: map[string]interface{}{
				"b": map[string]interface{}{
					"a": "bye",
				},
			},
		},
		{
			name: "c.b.a",
			args: args{
				s:     []string{"c", "b", "a"},
				value: "allo",
				dest:  map[string]interface{}{},
			},
			want: map[string]interface{}{
				"c": map[string]interface{}{
					"b": map[string]interface{}{
						"a": "allo",
					},
				},
			},
		},
		{
			name: "c.b.a",
			args: args{
				s:     []string{"a", "b", "c"},
				value: "bye",
				dest: map[string]interface{}{
					"a": map[string]interface{}{
						"b": map[string]interface{}{
							"d": "allo",
						},
					},
				},
			},
			want: map[string]interface{}{
				"a": map[string]interface{}{
					"b": map[string]interface{}{
						"d": "allo",
						"c": "bye",
					},
				},
			},
		},
		{
			name: "c.b. asdasa",
			args: args{
				s:     []string{"files", "fileb", "size"},
				value: "20",
				dest: map[string]interface{}{
					"app": map[string]interface{}{
						"source": map[string]interface{}{
							"git":     "github.com",
							"version": "1.2.3",
						},
						"attributes": map[string]interface{}{
							"name": "myapp",
						},
					},
					"files": map[string]interface{}{
						"filea": map[string]interface{}{
							"name": "filea",
							"size": "12",
						},
						"fileb": map[string]interface{}{
							"name": "fileb",
						},
						"filec": map[string]interface{}{
							"name": "filec",
						},
					},
				},
			},
			want: map[string]interface{}{
				"app": map[string]interface{}{
					"source": map[string]interface{}{
						"git":     "github.com",
						"version": "1.2.3",
					},
					"attributes": map[string]interface{}{
						"name": "myapp",
					},
				},
				"files": map[string]interface{}{
					"filea": map[string]interface{}{
						"name": "filea",
						"size": "12",
					},
					"fileb": map[string]interface{}{
						"name": "fileb",
						"size": "20",
					},
					"filec": map[string]interface{}{
						"name": "filec",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := insertInMapRecursion(tt.args.s, tt.args.value, tt.args.dest)
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("insertInMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
