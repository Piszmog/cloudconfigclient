package cloudconfigclient_test

import (
	"errors"
	"github.com/Piszmog/cloudconfigclient/v2"
	"github.com/stretchr/testify/require"
	"net/http"
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
			options: []cloudconfigclient.Option{cloudconfigclient.Local(&http.Client{}, []string{"http:localhost:8888"})},
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
	tests := []struct {
		name string
	}{
		// TODO: test cases
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
		})
	}
}
