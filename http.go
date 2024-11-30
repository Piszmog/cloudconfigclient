package cloudconfigclient

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"gopkg.in/yaml.v3"
)

// HTTPClient is a wrapper for http.Client.
type HTTPClient struct {
	*http.Client
	// BaseURL is the base URL for the Config Server.
	BaseURL string
	// Authorization is the authorization header value for the Config Server. If not provided, no authorization header is not explicitly set.
	// If the client is using OAuth2, the authorization header is set automatically.
	Authorization string
}

// ErrResourceNotFound is a special error that is used to propagate 404s.
var ErrResourceNotFound = errors.New("failed to find resource")

const (
	failedToDecodeMessage = "failed to decode response from url: %w"
)

// GetResource performs a http.MethodGet operation. Builds the URL based on the provided paths and params. Deserializes
// the response to the specified destination.
//
// Capable of unmarshalling YAML, JSON, and XML. If file type is of another type, use GetResourceRaw instead.
func (h *HTTPClient) GetResource(paths []string, params map[string]string, dest interface{}) error {
	if len(paths) == 0 {
		return errors.New("no resource specified to be retrieved")
	}
	resp, err := h.Get(paths, params)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return ErrResourceNotFound
	}
	if resp.StatusCode != http.StatusOK {
		var b []byte
		b, err = io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read body with status code '%d': %w", resp.StatusCode, err)
		}
		return fmt.Errorf("server responded with status code '%d' and body '%s'", resp.StatusCode, b)
	}
	if err = decodeResponseBody(paths[len(paths)-1], resp, dest); err != nil {
		return err
	}
	return nil
}

func decodeResponseBody(file string, resp *http.Response, dest interface{}) error {
	if strings.Contains(file, ".yml") || strings.Contains(file, ".yaml") {
		if err := yaml.NewDecoder(resp.Body).Decode(dest); err != nil {
			return fmt.Errorf(failedToDecodeMessage, err)
		}
	} else if strings.Contains(file, ".xml") {
		if err := xml.NewDecoder(resp.Body).Decode(dest); err != nil {
			return fmt.Errorf(failedToDecodeMessage, err)
		}
	} else {
		if err := json.NewDecoder(resp.Body).Decode(dest); err != nil {
			return fmt.Errorf(failedToDecodeMessage, err)
		}
	}
	return nil
}

// GetResourceRaw performs a http.MethodGet operation. Builds the URL based on the provided paths and params. Returns
// the byte slice response.
func (h *HTTPClient) GetResourceRaw(paths []string, params map[string]string) ([]byte, error) {
	if len(paths) == 0 {
		return nil, errors.New("no resource specified to be retrieved")
	}
	resp, err := h.Get(paths, params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrResourceNotFound
	}
	var b []byte
	b, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body with status code '%d': %w", resp.StatusCode, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server responded with status code '%d' and body '%s'", resp.StatusCode, b)
	}
	return b, nil
}

// Get performs a http.MethodGet operation. Builds the URL based on the provided paths and params.
func (h *HTTPClient) Get(paths []string, params map[string]string) (*http.Response, error) {
	fullURL, err := newURL(h.BaseURL, paths, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create url: %w", err)
	}
	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for %s: %w", fullURL, err)
	}
	if h.Authorization != "" {
		req.Header.Set("Authorization", h.Authorization)
	}
	response, err := h.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve from %s: %w", fullURL, err)
	}
	return response, nil
}

func newURL(baseURL string, paths []string, params map[string]string) (string, error) {
	parseURL, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse url %s: %w", baseURL, err)
	}
	if len(paths) > 0 {
		for _, p := range paths {
			parseURL.Path = path.Join(parseURL.Path, p)
		}
	}
	if params != nil {
		query := parseURL.Query()
		for key, value := range params {
			query.Set(key, value)
		}
		parseURL.RawQuery = query.Encode()
	}
	return parseURL.String(), nil
}
