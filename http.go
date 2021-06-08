package cloudconfigclient

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
)

// HTTPClient is a wrapper for http.Client.
type HTTPClient struct {
	*http.Client
	BaseURL string
}

var notFoundError = errors.New("failed to find resource")

// GetResource performs a http.MethodGet operation. Builds the URL based on the provided paths and params. Deserializes
// the response to the specified destination.
func (h *HTTPClient) GetResource(paths []string, params map[string]string, dest interface{}) error {
	resp, err := h.Get(paths, params)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return notFoundError
	}
	if resp.StatusCode != http.StatusOK {
		var b []byte
		b, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read body with status code '%d': %w", resp.StatusCode, err)
		}
		return fmt.Errorf("server responded with status code '%d' and body '%s'", resp.StatusCode, b)
	}
	if strings.Contains(paths[len(paths)-1], ".yml") || strings.Contains(paths[len(paths)-1], ".yaml") {
		if err = yaml.NewDecoder(resp.Body).Decode(dest); err != nil {
			return fmt.Errorf("failed to decode response from url: %w", err)
		}
	} else if strings.Contains(paths[len(paths)-1], ".xml") {
		if err = xml.NewDecoder(resp.Body).Decode(dest); err != nil {
			return fmt.Errorf("failed to decode response from url: %w", err)
		}
	} else {
		if err = json.NewDecoder(resp.Body).Decode(dest); err != nil {
			return fmt.Errorf("failed to decode response from url: %w", err)
		}
	}
	return nil
}

// GetResourceRaw performs a http.MethodGet operation. Builds the URL based on the provided paths and params. Returns
// the byte slice response.
func (h *HTTPClient) GetResourceRaw(paths []string, params map[string]string) ([]byte, error) {
	resp, err := h.Get(paths, params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, notFoundError
	}
	var b []byte
	b, err = ioutil.ReadAll(resp.Body)
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
	fullUrl, err := newURL(h.BaseURL, paths, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create url: %w", err)
	}
	response, err := h.Client.Get(fullUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve from %s: %w", fullUrl, err)
	}
	return response, nil
}

func newURL(baseURL string, paths []string, params map[string]string) (string, error) {
	parseUrl, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse url %s: %w", baseURL, err)
	}
	if paths != nil {
		for _, p := range paths {
			parseUrl.Path = path.Join(parseUrl.Path, p)
		}
	}
	if params != nil {
		query := parseUrl.Query()
		for key, value := range params {
			query.Set(key, value)
		}
		parseUrl.RawQuery = query.Encode()
	}
	return parseUrl.String(), nil
}
