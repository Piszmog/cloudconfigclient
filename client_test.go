package cloudconfigclient_test

import (
	"context"
	"errors"
	"github.com/Piszmog/cloudconfigclient/v2"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2/clientcredentials"
	"net/http"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		options []cloudconfigclient.Option
		err     error
	}{
		{
			name:    "No Options",
			options: nil,
			err:     errors.New("at least one option must be provided"),
		},
		{
			name:    "Option Error",
			options: []cloudconfigclient.Option{cloudconfigclient.LocalEnv(&http.Client{})},
			err:     errors.New("no local Config Server URLs provided in environment variable CONFIG_SERVER_URLS"),
		},
		{
			name:    "Created",
			options: []cloudconfigclient.Option{cloudconfigclient.Local(&http.Client{}, "http:localhost:8888")},
			err:     nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client, err := cloudconfigclient.New(test.options...)
			if err != nil {
				require.Equal(t, test.err, err)
				require.Nil(t, client)
			} else {
				require.NoError(t, err)
				require.NotNil(t, client)
			}
		})
	}
}

func TestOption(t *testing.T) {
	oauthConfig := clientcredentials.Config{ClientID: "clientId", ClientSecret: "secret", TokenURL: "http://token"}
	oauthClient := oauthConfig.Client(context.Background())
	tests := []struct {
		name     string
		setup    func()
		cleanup  func()
		option   cloudconfigclient.Option
		expected []*cloudconfigclient.HTTPClient
		err      error
	}{
		{
			name: "LocalEnv",
			setup: func() {
				os.Setenv(cloudconfigclient.EnvironmentLocalConfigServerUrls, "http://localhost:8880")
			},
			cleanup: func() {
				os.Unsetenv(cloudconfigclient.EnvironmentLocalConfigServerUrls)
			},
			option:   cloudconfigclient.LocalEnv(&http.Client{}),
			expected: []*cloudconfigclient.HTTPClient{{BaseURL: "http://localhost:8880", Client: &http.Client{}}},
		},
		{
			name: "LocalEnv Multiple",
			setup: func() {
				os.Setenv(cloudconfigclient.EnvironmentLocalConfigServerUrls, "http://localhost:8880,http://localhost:8888")
			},
			cleanup: func() {
				os.Unsetenv(cloudconfigclient.EnvironmentLocalConfigServerUrls)
			},
			option: cloudconfigclient.LocalEnv(&http.Client{}),
			expected: []*cloudconfigclient.HTTPClient{
				{BaseURL: "http://localhost:8880", Client: &http.Client{}},
				{BaseURL: "http://localhost:8888", Client: &http.Client{}},
			},
		},
		{
			name:   "LocalEnv Error",
			option: cloudconfigclient.LocalEnv(&http.Client{}),
			err:    errors.New("no local Config Server URLs provided in environment variable CONFIG_SERVER_URLS"),
		},
		{
			name:     "Local",
			option:   cloudconfigclient.Local(&http.Client{}, "http://localhost:8880"),
			expected: []*cloudconfigclient.HTTPClient{{BaseURL: "http://localhost:8880", Client: &http.Client{}}},
		},
		{
			name:   "Local Multiple",
			option: cloudconfigclient.Local(&http.Client{}, "http://localhost:8880", "http://localhost:8888"),
			expected: []*cloudconfigclient.HTTPClient{
				{BaseURL: "http://localhost:8880", Client: &http.Client{}},
				{BaseURL: "http://localhost:8888", Client: &http.Client{}},
			},
		},
		{
			name: "DefaultCFService",
			setup: func() {
				os.Setenv("VCAP_SERVICES", `{
      "p.config-server": [
        {
          "credentials": {
            "uri": "http://config",
            "client_secret": "secret",
            "client_id": "clientId",
            "access_token_uri": "http://token"
          }
        }
      ]
    }`)
			},
			cleanup: func() {
				os.Unsetenv("VCAP_SERVICES")
			},
			option: cloudconfigclient.DefaultCFService(),
			expected: []*cloudconfigclient.HTTPClient{{
				BaseURL: "http://config",
				Client:  oauthClient,
			}},
		},
		{
			name: "DefaultCFService Old Service",
			setup: func() {
				os.Setenv("VCAP_SERVICES", `{
      "p-config-server": [
        {
          "credentials": {
            "uri": "http://config",
            "client_secret": "secret",
            "client_id": "clientId",
            "access_token_uri": "http://token"
          }
        }
      ]
    }`)
			},
			cleanup: func() {
				os.Unsetenv("VCAP_SERVICES")
			},
			option: cloudconfigclient.DefaultCFService(),
			expected: []*cloudconfigclient.HTTPClient{{
				BaseURL: "http://config",
				Client:  oauthClient,
			}},
		},
		{
			name: "DefaultCFService Not Found",
			setup: func() {
				os.Setenv("VCAP_SERVICES", `{}`)
			},
			cleanup: func() {
				os.Unsetenv("VCAP_SERVICES")
			},
			option: cloudconfigclient.DefaultCFService(),
			err:    errors.New("neither p-config-server or p.config-server exist in environment variable 'VCAP_SERVICES'"),
		},
		{
			name:   "DefaultCFService Error",
			option: cloudconfigclient.DefaultCFService(),
			err:    errors.New("failed to parse 'VCAP_SERVICES': failed to unmarshal JSON: unexpected end of JSON input"),
		},
		{
			name: "CFService",
			setup: func() {
				os.Setenv("VCAP_SERVICES", `{
      "config-server": [
        {
          "credentials": {
            "uri": "http://config",
            "client_secret": "secret",
            "client_id": "clientId",
            "access_token_uri": "http://token"
          }
        }
      ]
    }`)
			},
			cleanup: func() {
				os.Unsetenv("VCAP_SERVICES")
			},
			option: cloudconfigclient.CFService("config-server"),
			expected: []*cloudconfigclient.HTTPClient{{
				BaseURL: "http://config",
				Client:  oauthClient,
			}},
		},
		{
			name:   "CFService Error",
			option: cloudconfigclient.CFService("config-server"),
			err:    errors.New("failed to parse 'VCAP_SERVICES': failed to unmarshal JSON: unexpected end of JSON input"),
		},
		{
			name: "CFService Not Found",
			setup: func() {
				os.Setenv("VCAP_SERVICES", `{
      "something-else": [
        {
          "credentials": {
            "uri": "http://config",
            "client_secret": "secret",
            "client_id": "clientId",
            "access_token_uri": "http://token"
          }
        }
      ]
    }`)
			},
			cleanup: func() {
				os.Unsetenv("VCAP_SERVICES")
			},
			option: cloudconfigclient.CFService("config-server"),
			err:    errors.New("failed to create cloud Client: service does not exist"),
		},
		{
			name:   "OAuth2",
			option: cloudconfigclient.OAuth2("http://config", "clientId", "secret", "http://token"),
			expected: []*cloudconfigclient.HTTPClient{{
				BaseURL: "http://config",
				Client:  oauthClient,
			}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.setup != nil {
				test.setup()
			}
			if test.cleanup != nil {
				defer test.cleanup()
			}
			var clients []*cloudconfigclient.HTTPClient
			err := test.option(&clients)
			if err != nil {
				require.Error(t, err)
				require.Equal(t, test.err.Error(), err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expected, clients)
			}
		})
	}
}
