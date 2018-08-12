package net

import (
	"github.com/Piszmog/cfservices"
	"testing"
)

func TestCreateOAuth2Client(t *testing.T) {
	credentials := &cfservices.Credentials{
		AccessTokenUri: "tokenUri",
		ClientSecret:   "clientSecret",
		ClientId:       "clientId",
	}
	client, err := CreateOAuth2Client(credentials)
	if err != nil {
		t.Errorf("failed to create oauth2 client with error %v", err)
	}
	if client == nil {
		t.Error("no oauth2 client returned")
	}
}

func TestCreateOAuth2ClientWhenCredentialsAreNil(t *testing.T) {
	client, err := CreateOAuth2Client(nil)
	if err == nil {
		t.Error("expected an error when no credentials are passed")
	}
	if client != nil {
		t.Error("able to create an oauth2 client with nil credentials")
	}
}

func TestCreateOauth2Config(t *testing.T) {
	credentials := &cfservices.Credentials{
		AccessTokenUri: "tokenUri",
		ClientSecret:   "clientSecret",
		ClientId:       "clientId",
	}
	config, err := CreateOAuth2Config(credentials)
	if err != nil {
		t.Errorf("failed to create oauth2 with errpr %v", err)
	}
	if config == nil {
		t.Error("failed to create oauth2 config")
	}
}

func TestCreateOauth2ConfigWhenCredentialsNil(t *testing.T) {
	config, err := CreateOAuth2Config(nil)
	if err == nil {
		t.Error("expected an error when passing nil credentials when creating oauth2 config")
	}
	if config != nil {
		t.Error("is able to create oauth2 config when credentials are nil")
	}
}
