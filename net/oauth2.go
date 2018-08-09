package net

import (
	"github.com/Piszmog/cfservices/credentials"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/clientcredentials"
	"net/http"
)

func CreateOAuth2Client(cred credentials.Credentials) *http.Client {
	config := CreateOauth2Config(&cred)
	return config.Client(context.Background())
}

func CreateOauth2Config(cred *credentials.Credentials) *clientcredentials.Config {
	return &clientcredentials.Config{
		ClientID:     cred.ClientId,
		ClientSecret: cred.ClientSecret,
		TokenURL:     cred.AccessTokenUri,
	}
}
