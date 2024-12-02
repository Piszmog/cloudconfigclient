package cloudconfigclient_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/ngaggi73/cloudconfigclient/v2"
	"github.com/stretchr/testify/require"
)

const (
	testJSONFile = `{
  "example":{
    "field":"value"
  }
}`
)

type file struct {
	Example example `json:"example"`
}

type example struct {
	Field string `json:"field"`
}

func TestClient_GetFile(t *testing.T) {
	tests := []struct {
		name     string
		checker  func(*testing.T, *http.Request)
		response *http.Response
		expected file
		err      error
	}{
		{
			name: "JSON File",
			checker: func(t *testing.T, request *http.Request) {
				require.Equal(t, "http://localhost:8888/default/default/directory/file.json?useDefaultLabel=true", request.URL.String())
			},
			response: NewMockHttpResponse(http.StatusOK, testJSONFile),
			expected: file{Example: example{Field: "value"}},
		},
		{
			name:     "Not Found",
			response: NewMockHttpResponse(http.StatusNotFound, ""),
			err:      errors.New("failed to find file in the Config Server"),
		},
		{
			name:     "Server Error",
			response: NewMockHttpResponse(http.StatusInternalServerError, ""),
			err:      errors.New("server responded with status code '500' and body ''"),
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

			var actual file
			err = client.GetFile("directory", "file.json", &actual)
			if err != nil {
				require.Error(t, err)
				require.Equal(t, test.err.Error(), err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expected, actual)
			}
		})
	}
}

func TestClient_GetFileFromBranch(t *testing.T) {
	tests := []struct {
		name     string
		checker  func(*testing.T, *http.Request)
		response *http.Response
		expected file
		err      error
	}{
		{
			name: "JSON File",
			checker: func(t *testing.T, request *http.Request) {
				require.Equal(t, "http://localhost:8888/default/default/branch/directory/file.json", request.URL.String())
			},
			response: NewMockHttpResponse(http.StatusOK, testJSONFile),
			expected: file{Example: example{Field: "value"}},
		},
		{
			name:     "Not Found",
			response: NewMockHttpResponse(http.StatusNotFound, ""),
			err:      errors.New("failed to find file in the Config Server"),
		},
		{
			name:     "Server Error",
			response: NewMockHttpResponse(http.StatusInternalServerError, ""),
			err:      errors.New("server responded with status code '500' and body ''"),
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

			var actual file
			err = client.GetFileFromBranch("branch", "directory", "file.json", &actual)
			if err != nil {
				require.Error(t, err)
				require.Equal(t, test.err.Error(), err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expected, actual)
			}
		})
	}
}

func TestClient_GetFileRaw(t *testing.T) {
	tests := []struct {
		name     string
		checker  func(*testing.T, *http.Request)
		response *http.Response
		expected []byte
		err      error
	}{
		{
			name: "Text File",
			checker: func(t *testing.T, request *http.Request) {
				require.Equal(t, "http://localhost:8888/default/default/directory/file.txt?useDefaultLabel=true", request.URL.String())
			},
			response: NewMockHttpResponse(http.StatusOK, "hello world"),
			expected: []byte("hello world"),
		},
		{
			name:     "Not Found",
			response: NewMockHttpResponse(http.StatusNotFound, ""),
			err:      errors.New("failed to find file in the Config Server"),
		},
		{
			name:     "Server Error",
			response: NewMockHttpResponse(http.StatusInternalServerError, ""),
			err:      errors.New("server responded with status code '500' and body ''"),
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

			actual, err := client.GetFileRaw("directory", "file.txt")
			if err != nil {
				require.Error(t, err)
				require.Equal(t, test.err.Error(), err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expected, actual)
			}
		})
	}
}

func TestClient_GetFileFromBranchRaw(t *testing.T) {
	tests := []struct {
		name     string
		checker  func(*testing.T, *http.Request)
		response *http.Response
		expected []byte
		err      error
	}{
		{
			name: "Text File",
			checker: func(t *testing.T, request *http.Request) {
				require.Equal(t, "http://localhost:8888/default/default/branch/directory/file.txt", request.URL.String())
			},
			response: NewMockHttpResponse(http.StatusOK, "hello world"),
			expected: []byte("hello world"),
		},
		{
			name:     "Not Found",
			response: NewMockHttpResponse(http.StatusNotFound, ""),
			err:      errors.New("failed to find file in the Config Server"),
		},
		{
			name:     "Server Error",
			response: NewMockHttpResponse(http.StatusInternalServerError, ""),
			err:      errors.New("server responded with status code '500' and body ''"),
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

			actual, err := client.GetFileFromBranchRaw("branch", "directory", "file.txt")
			if err != nil {
				require.Error(t, err)
				require.Equal(t, test.err.Error(), err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expected, actual)
			}
		})
	}
}
