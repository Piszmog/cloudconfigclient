package configuration

import (
	"testing"
)

func TestCreateClient(t *testing.T) {
	client := CreateClient("http://localhost:8080")
	if client == nil && client.BaseUrls[0] == "http://localhost:8080" {
		t.Errorf("failed to create configuration client")
	}
}
