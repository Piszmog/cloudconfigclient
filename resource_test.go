package cloudconfigclient

import (
	"errors"
	"testing"
)

const (
	testJSONFile = `{
  "example":{
    "field":"value"
  }
}`
)

type file struct {
	Example example `json:"example"`
}

type example struct {
	Field string `json:"field"`
}

func TestConfigClient_GetFile(t *testing.T) {
	configClient := createMockConfigClient(200, testJSONFile, nil)
	var file file
	err := configClient.GetFile("directory", "file.json", &file)
	if err != nil {
		t.Errorf("failed to retrieve file with error %v", err)
	}
	if file.Example.Field != "value" {
		t.Error("failed to retrieve file")
	}
}

func TestConfigClient_GetFileWhen404(t *testing.T) {
	configClient := createMockConfigClient(404, "", nil)
	var file file
	err := configClient.GetFile("directory", "file.json", &file)
	if err == nil {
		t.Error("expected an error to occur")
	}
	if file.Example.Field == "value" {
		t.Error("retrieved configuration when not found")
	}
}

func TestConfigClient_GetFileWhenError(t *testing.T) {
	configClient := createMockConfigClient(500, "", errors.New("failed"))
	var file file
	err := configClient.GetFile("directory", "file.json", &file)
	if err == nil {
		t.Error("expected an error to occur")
	}
	if file.Example.Field == "value" {
		t.Error("retrieved configuration when not found")
	}
}

func TestConfigClient_GetFileWhenNoErrorBut500(t *testing.T) {
	configClient := createMockConfigClient(500, "", nil)
	var file file
	err := configClient.GetFile("directory", "file.json", &file)
	if err == nil {
		t.Error("expected an error to occur")
	}
	if file.Example.Field == "value" {
		t.Error("retrieved configuration when not found")
	}
}

func TestConfigClient_GetFileInvalidResponseBody(t *testing.T) {
	configClient := createMockConfigClient(200, "", nil)
	var file file
	err := configClient.GetFile("directory", "file.json", &file)
	if err == nil {
		t.Error("expected an error to occur")
	}
	if file.Example.Field == "value" {
		t.Error("retrieved configuration when not found")
	}
}

func TestConfigClient_GetFileFromBranch(t *testing.T) {
	configClient := createMockConfigClient(200, testJSONFile, nil)
	var file file
	err := configClient.GetFileFromBranch("branch", "directory", "file.json", &file)
	if err != nil {
		t.Errorf("failed to retrieve file with error %v", err)
	}
	if file.Example.Field != "value" {
		t.Error("failed to retrieve file")
	}
}

func TestConfigClient_GetFileFromBranchWhen404(t *testing.T) {
	configClient := createMockConfigClient(404, "", nil)
	var file file
	err := configClient.GetFileFromBranch("branch", "directory", "file.json", &file)
	if err == nil {
		t.Error("expected an error to occur")
	}
	if file.Example.Field == "value" {
		t.Error("retrieved configuration when not found")
	}
}

func TestConfigClient_GetFileFromBranchWhenError(t *testing.T) {
	configClient := createMockConfigClient(500, "", errors.New("failed"))
	var file file
	err := configClient.GetFileFromBranch("branch", "directory", "file.json", &file)
	if err == nil {
		t.Error("expected an error to occur")
	}
	if file.Example.Field == "value" {
		t.Error("retrieved configuration when not found")
	}
}

func TestConfigClient_GetFileFromBranchWhenNoErrorBut500(t *testing.T) {
	configClient := createMockConfigClient(500, "", nil)
	var file file
	err := configClient.GetFileFromBranch("branch", "directory", "file.json", &file)
	if err == nil {
		t.Error("expected an error to occur")
	}
	if file.Example.Field == "value" {
		t.Error("retrieved configuration when not found")
	}
}

func TestConfigClient_GetFileFromBranchInvalidResponseBody(t *testing.T) {
	configClient := createMockConfigClient(200, "", nil)
	var file file
	err := configClient.GetFileFromBranch("branch", "directory", "file.json", &file)
	if err == nil {
		t.Error("expected an error to occur")
	}
	if file.Example.Field == "value" {
		t.Error("retrieved configuration when not found")
	}
}
