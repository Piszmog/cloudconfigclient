package client

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func CreateMockClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func TestClient_Get(t *testing.T) {
	httpClient := CreateMockClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			// Send response to be tested
			Body: ioutil.NopCloser(bytes.NewBufferString(`OK`)),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})
	client := &Client{
		configUri:  "http://localhost:8080",
		httpClient: httpClient,
	}
	resp, err := client.Get("some", "path")
	if err != nil {
		t.Errorf("failed to call the mock server with error %v", err)
	}
	if resp == nil {
		t.Errorf("expected a response body")
	}
	defer resp.Body.Close()
	byteBody, _ := ioutil.ReadAll(resp.Body)
	if string(byteBody) != "OK" {
		t.Error("failed to read body")
	}
}
