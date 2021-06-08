package cloudconfigclient_test

import (
	"github.com/Piszmog/cloudconfigclient/v2"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestClient_Get(t *testing.T) {
	httpClient := NewMockHttpClient(func(req *http.Request) *http.Response {
		assert.Equal(t, "/some/path", req.URL.RequestURI())
		return NewMockHttpResponse(200, "OK")
	})
	client := &cloudconfigclient.Client{
		ConfigUri:  "http://localhost:8080",
		HttpClient: httpClient,
	}
	resp, err := client.Get("some", "path")
	assert.NoError(t, err, "failed to call the mock server with error")
	assert.NotNil(t, resp)
	defer resp.Body.Close()
	byteBody, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, "OK", string(byteBody))
}

func TestClient_Get_QueryParam(t *testing.T) {
	httpClient := NewMockHttpClient(func(req *http.Request) *http.Response {
		assert.Equal(t, "/some/path/file.txt?key=value", req.URL.RequestURI())
		return NewMockHttpResponse(200, "OK")
	})
	client := &cloudconfigclient.Client{
		ConfigUri:  "http://localhost:8080",
		HttpClient: httpClient,
	}
	resp, err := client.Get("some", "path", "file.txt?key=value")
	assert.NoError(t, err, "failed to call the mock server with error")
	assert.NotNil(t, resp)
	defer resp.Body.Close()
	byteBody, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, "OK", string(byteBody))
}
