package cloudconfigclient_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
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
	cloudClient := new(mockCloudClient)
	response := NewMockHttpResponse(200, testJSONFile)
	cloudClient.On("Get", []string{"default", "default", "directory", "file.json?useDefaultLabel=true"}).Return(response, nil)
	client := NewConfigClient(cloudClient)
	var f file
	err := client.GetFile("directory", "file.json", &f)
	assert.NoError(t, err)
	assert.Equal(t, "value", f.Example.Field)
}

func TestConfigClient_GetFileWhen404(t *testing.T) {
	cloudClient := new(mockCloudClient)
	response := NewMockHttpResponse(404, "")
	cloudClient.On("Get", []string{"default", "default", "directory", "file.json?useDefaultLabel=true"}).
		Return(response, nil)
	client := NewConfigClient(cloudClient)
	var f file
	err := client.GetFile("directory", "file.json", &f)
	assert.Error(t, err)
	assert.Empty(t, f.Example.Field)
}

func TestConfigClient_GetFileWhenError(t *testing.T) {
	cloudClient := new(mockCloudClient)
	response := NewMockHttpResponse(500, "")
	cloudClient.On("Get", []string{"default", "default", "directory", "file.json?useDefaultLabel=true"}).
		Return(response, errors.New("failed"))
	client := NewConfigClient(cloudClient)
	var f file
	err := client.GetFile("directory", "file.json", &f)
	assert.Error(t, err)
	assert.Empty(t, f.Example.Field)
}

func TestConfigClient_GetFileWhenNoErrorBut500(t *testing.T) {
	cloudClient := new(mockCloudClient)
	response := NewMockHttpResponse(500, "")
	cloudClient.On("Get", []string{"default", "default", "directory", "file.json?useDefaultLabel=true"}).
		Return(response, nil)
	client := NewConfigClient(cloudClient)
	var file file
	err := client.GetFile("directory", "file.json", &file)
	assert.Error(t, err)
	assert.Empty(t, file.Example.Field)
}

func TestConfigClient_GetFileInvalidResponseBody(t *testing.T) {
	cloudClient := new(mockCloudClient)
	response := NewMockHttpResponse(200, "")
	cloudClient.On("Get", []string{"default", "default", "directory", "file.json?useDefaultLabel=true"}).
		Return(response, nil)
	client := NewConfigClient(cloudClient)
	var file file
	err := client.GetFile("directory", "file.json", &file)
	assert.Error(t, err)
	assert.Empty(t, file.Example.Field)
}

func TestConfigClient_GetFileFromBranch(t *testing.T) {
	cloudClient := new(mockCloudClient)
	response := NewMockHttpResponse(200, testJSONFile)
	cloudClient.On("Get", []string{"default", "default", "branch", "directory", "file.json"}).
		Return(response, nil)
	client := NewConfigClient(cloudClient)
	var f file
	err := client.GetFileFromBranch("branch", "directory", "file.json", &f)
	assert.NoError(t, err)
	assert.Equal(t, "value", f.Example.Field)
}

func TestConfigClient_GetFileFromBranchWhen404(t *testing.T) {
	cloudClient := new(mockCloudClient)
	response := NewMockHttpResponse(404, "")
	cloudClient.On("Get", []string{"default", "default", "branch", "directory", "file.json"}).
		Return(response, nil)
	client := NewConfigClient(cloudClient)
	var f file
	err := client.GetFileFromBranch("branch", "directory", "file.json", &f)
	assert.Error(t, err)
	assert.Empty(t, f.Example.Field)
}

func TestConfigClient_GetFileFromBranchWhenError(t *testing.T) {
	cloudClient := new(mockCloudClient)
	response := NewMockHttpResponse(500, "")
	cloudClient.On("Get", []string{"default", "default", "branch", "directory", "file.json"}).
		Return(response, errors.New("failed"))
	client := NewConfigClient(cloudClient)
	var f file
	err := client.GetFileFromBranch("branch", "directory", "file.json", &f)
	assert.Error(t, err)
	assert.Empty(t, f.Example.Field)
}

func TestConfigClient_GetFileFromBranchWhenNoErrorBut500(t *testing.T) {
	cloudClient := new(mockCloudClient)
	response := NewMockHttpResponse(500, "")
	cloudClient.On("Get", []string{"default", "default", "branch", "directory", "file.json"}).
		Return(response, nil)
	client := NewConfigClient(cloudClient)
	var f file
	err := client.GetFileFromBranch("branch", "directory", "file.json", &f)
	assert.Error(t, err)
	assert.Empty(t, f.Example.Field)
}

func TestConfigClient_GetFileFromBranchInvalidResponseBody(t *testing.T) {
	cloudClient := new(mockCloudClient)
	response := NewMockHttpResponse(200, "")
	cloudClient.On("Get", []string{"default", "default", "branch", "directory", "file.json"}).
		Return(response, nil)
	client := NewConfigClient(cloudClient)
	var f file
	err := client.GetFileFromBranch("branch", "directory", "file.json", &f)
	assert.Error(t, err)
	assert.Empty(t, f.Example.Field)
}
