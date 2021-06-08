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
				require.Equal(t, test.err.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
