package net

import (
	"net"
	"net/http"
	"strings"
	"time"
)

func CreateUrl(baseUrl string, uriVariables ...string) string {
	url := strings.TrimRight(baseUrl, "/")
	for _, uriVariable := range uriVariables {
		url = url + "/" + uriVariable
	}
	return url
}

func JoinProfiles(profiles []string) string {
	return strings.Join(profiles, ",")
}

func CreateDefaultHttpClient() *http.Client {
	return CreateHttpClient(5*time.Second, 30*time.Second, 5*time.Second, 90*time.Second)
}

func CreateHttpClient(timeout time.Duration, keepAlive time.Duration, tlsHandshakeTimeout time.Duration, idleConnection time.Duration) *http.Client {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   timeout,
			KeepAlive: keepAlive,
			DualStack: true,
		}).DialContext,
		TLSHandshakeTimeout: tlsHandshakeTimeout,
		IdleConnTimeout:     idleConnection,
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
	}
	return &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}
}
