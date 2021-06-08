package cloudconfigclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type httpClient struct {
	baseURL string
	client  *http.Client
}

var notFoundError = errors.New("failed to find resource")

func (h *httpClient) getResource(paths []string, params map[string]string, dest interface{}) error {
	resp, err := h.get(paths, params)
	if err != nil {
		return err
	}
	defer closeResource(resp.Body)
	if resp.StatusCode == http.StatusNotFound {
		return notFoundError
	}
	if resp.StatusCode != http.StatusOK {
		var b []byte
		b, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read body with status code %d: %w", resp.StatusCode, err)
		}
		return fmt.Errorf("server responded with status code %d and body %s", resp.StatusCode, b)
	}
	if strings.Contains(paths[len(paths)-1], ".yml") || strings.Contains(paths[len(paths)-1], ".yaml") {
		if err = yaml.NewDecoder(resp.Body).Decode(dest); err != nil {
			return fmt.Errorf("failed to decode response from url: %w", err)
		}
	} else {
		if err = json.NewDecoder(resp.Body).Decode(dest); err != nil {
			return fmt.Errorf("failed to decode response from url: %w", err)
		}
	}
	return nil
}

func (h *httpClient) get(paths []string, params map[string]string) (*http.Response, error) {
	fullUrl, err := newURL(h.baseURL, paths, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create url: %w", err)
	}
	response, err := h.client.Get(fullUrl)
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
	for _, p := range paths {
		parseUrl.Path = path.Join(parseUrl.Path, p)
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
