package cloudconfigclient_test

import (
	"bytes"
	"errors"
	"github.com/Piszmog/cloudconfigclient/v2"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"testing"
)

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// NewMockHttpClient creates a mocked HTTP Client
func NewMockHttpClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}

// NewMockHttpResponse creates a mocked HTTP response
func NewMockHttpResponse(code int, body string) *http.Response {
	return &http.Response{
		StatusCode: code,
		// Send response to be tested
		Body: ioutil.NopCloser(bytes.NewBufferString(body)),
		// Must be set to non-nil value or it panics
		Header: make(http.Header),
	}
}

func TestHTTPClient_Get(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		paths    []string
		params   map[string]string
		response *http.Response
		checker  func(*testing.T, *http.Request)
		err      error
	}{
		{
			name:    "Invalid URL",
			baseURL: "\n",
			err:     errors.New("failed to create url: failed to parse url \n: parse \"\\n\": net/url: invalid control character in URL"),
		},
		{
			name:    "HTTP Error",
			baseURL: "http://foobar",
			err:     errors.New("failed to retrieve from http://foobar: Get \"http://foobar\": http: RoundTripper implementation (cloudconfigclient_test.RoundTripFunc) returned a nil *Response with a nil error"),
		},
		{
			name:     "Correct Request URI",
			baseURL:  "http://something",
			response: NewMockHttpResponse(200, ""),
			checker: func(t *testing.T, request *http.Request) {
				require.Equal(t, "/", request.URL.RequestURI())
			},
		},
		{
			name:     "Correct Request URI With Path",
			baseURL:  "http://something",
			paths:    []string{"/foo", "/bar"},
			response: NewMockHttpResponse(200, ""),
			checker: func(t *testing.T, request *http.Request) {
				require.Equal(t, "/foo/bar", request.URL.RequestURI())
			},
		},
		{
			name:     "Correct Request URI With Params",
			baseURL:  "http://something",
			params:   map[string]string{"field": "value"},
			response: NewMockHttpResponse(200, ""),
			checker: func(t *testing.T, request *http.Request) {
				require.Equal(t, "/?field=value", request.URL.RequestURI())
			},
		},
		{
			name:     "Correct Request URI With Path and Params",
			baseURL:  "http://something",
			paths:    []string{"/foo", "/bar"},
			params:   map[string]string{"field": "value"},
			response: NewMockHttpResponse(200, ""),
			checker: func(t *testing.T, request *http.Request) {
				require.Equal(t, "/foo/bar?field=value", request.URL.RequestURI())
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := NewMockHttpClient(func(req *http.Request) *http.Response {
				if test.checker != nil {
					test.checker(t, req)
				}
				return test.response
			})
			httpClient := cloudconfigclient.HTTPClient{BaseURL: test.baseURL, Client: client}
			_, err := httpClient.Get(test.paths, test.params)
			if err != nil {
				require.Error(t, err)
				require.Equal(t, test.err.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestHTTPClient_GetResource(t *testing.T) {
	tests := []struct {
		name        string
		paths       []string
		params      map[string]string
		destination interface{}
		response    *http.Response
		expected    interface{}
		err         error
	}{
		{
			name: "HTTP Error",
			err:  errors.New("failed to retrieve from http://something: Get \"http://something\": http: RoundTripper implementation (cloudconfigclient_test.RoundTripFunc) returned a nil *Response with a nil error"),
		},
		{
			name:     "Internal Server Error",
			response: NewMockHttpResponse(http.StatusInternalServerError, "Invalid HTTP Call"),
			err:      errors.New("server responded with status code '500' and body 'Invalid HTTP Call'"),
		},
		{
			name: "HTTP Error",
			err:  errors.New("failed to retrieve from http://something: Get \"http://something\": http: RoundTripper implementation (cloudconfigclient_test.RoundTripFunc) returned a nil *Response with a nil error"),
		},
		{
			name:        "YAML Response",
			paths:       []string{"file.yaml"},
			params:      map[string]string{"useDefault": "true"},
			destination: make(map[string]interface{}),
			response:    NewMockHttpResponse(http.StatusOK, `foo: bar`),
			expected:    map[string]interface{}{"foo": "bar"},
		},
		{
			name:        "YAML Response Malformed",
			paths:       []string{"file.yaml"},
			params:      map[string]string{"useDefault": "true"},
			destination: make(map[string]interface{}),
			response:    NewMockHttpResponse(http.StatusOK, ""),
			err:         errors.New("failed to decode response from url: EOF"),
		},
		{
			name:        "YML Response",
			paths:       []string{"file.yml"},
			params:      map[string]string{"useDefault": "true"},
			destination: make(map[string]interface{}),
			response:    NewMockHttpResponse(http.StatusOK, `foo: bar`),
			expected:    map[string]interface{}{"foo": "bar"},
		},
		{
			name:        "JSON Response",
			paths:       []string{"file.json"},
			params:      map[string]string{"useDefault": "true"},
			destination: make(map[string]interface{}),
			response:    NewMockHttpResponse(http.StatusOK, `{"foo":"bar"}`),
			expected:    map[string]interface{}{"foo": "bar"},
		},
		{
			name:        "JSON Response Malformed",
			paths:       []string{"file.json"},
			params:      map[string]string{"useDefault": "true"},
			destination: make(map[string]interface{}),
			response:    NewMockHttpResponse(http.StatusOK, `{"foo":"bar"`),
			err:         errors.New("failed to decode response from url: unexpected EOF"),
		},
		{
			name:        "XML Response",
			paths:       []string{"file.xml"},
			params:      map[string]string{"useDefault": "true"},
			destination: new(xmlResp),
			response:    NewMockHttpResponse(http.StatusOK, `"<data><foo>bar</foo></data>"`),
			expected:    &xmlResp{Foo: "bar"},
		},
		{
			name:        "XML Response Malformed",
			paths:       []string{"file.xml"},
			params:      map[string]string{"useDefault": "true"},
			destination: new(xmlResp),
			response:    NewMockHttpResponse(http.StatusOK, `"<data><foo>bar</foo></data"`),
			err:         errors.New("failed to decode response from url: XML syntax error on line 1: invalid characters between </data and >"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := NewMockHttpClient(func(req *http.Request) *http.Response {
				return test.response
			})
			httpClient := cloudconfigclient.HTTPClient{BaseURL: "http://something", Client: client}
			err := httpClient.GetResource(test.paths, test.params, &test.destination)
			if err != nil {
				require.Error(t, err)
				require.Equal(t, test.err.Error(), err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expected, test.destination)
			}
		})
	}
}

func TestHTTPClient_GetResourceRaw(t *testing.T) {
	tests := []struct {
		name     string
		paths    []string
		params   map[string]string
		response *http.Response
		expected []byte
		err      error
	}{
		{
			name: "HTTP Error",
			err:  errors.New("failed to retrieve from http://something: Get \"http://something\": http: RoundTripper implementation (cloudconfigclient_test.RoundTripFunc) returned a nil *Response with a nil error"),
		},
		{
			name:     "Internal Server Error",
			response: NewMockHttpResponse(http.StatusInternalServerError, "Invalid HTTP Call"),
			err:      errors.New("server responded with status code '500' and body 'Invalid HTTP Call'"),
		},
		{
			name: "HTTP Error",
			err:  errors.New("failed to retrieve from http://something: Get \"http://something\": http: RoundTripper implementation (cloudconfigclient_test.RoundTripFunc) returned a nil *Response with a nil error"),
		},
		{
			name:     "Text Response",
			paths:    []string{"file.text"},
			params:   map[string]string{"useDefault": "true"},
			response: NewMockHttpResponse(http.StatusOK, `foo-bar`),
			expected: []byte("foo-bar"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := NewMockHttpClient(func(req *http.Request) *http.Response {
				return test.response
			})
			httpClient := cloudconfigclient.HTTPClient{BaseURL: "http://something", Client: client}
			resp, err := httpClient.GetResourceRaw(test.paths, test.params)
			if err != nil {
				require.Error(t, err)
				require.Equal(t, test.err.Error(), err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expected, resp)
			}
		})
	}
}

type xmlResp struct {
	Foo string `xml:"foo"`
}
