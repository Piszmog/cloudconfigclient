# Go Config Server Client
[![Build Status](https://travis-ci.org/Piszmog/cloudconfigclient.svg?branch=develop)](https://travis-ci.org/Piszmog/cloudconfigclient)
[![Coverage Status](https://coveralls.io/repos/github/Piszmog/cloudconfigclient/badge.svg?branch=develop)](https://coveralls.io/github/Piszmog/cloudconfigclient?branch=develop)
[![Go Report Card](https://goreportcard.com/badge/github.com/Piszmog/cloudconfigclient)](https://goreportcard.com/report/github.com/Piszmog/cloudconfigclient)
[![GitHub release](https://img.shields.io/github/release/Piszmog/cloudconfigclient.svg)](https://github.com/Piszmog/cloudconfigclient/releases/latest)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Go library for Spring Config Server. Inspired by the Java library [Cloud Config Client](https://github.com/Piszmog/cloud-config-client)

## Description
Spring's Config Server provides way to externalize configurations of applications. Spring's
[Spring Cloud Config Client](https://github.com/spring-cloud/spring-cloud-config/tree/master/spring-cloud-config-client)
can be used to load the base configurations an application requires to function.

This library provides clients the ability to load Configurations and Files from the Config Server

## Example Usage
Below is an example usage of the library to retrieve a file from the Config Server and to retrieve the application's configurations

* For local config client, ensure `CONFIG_SERVER_URLS` is set
  * `CONFIG_SERVER_URLS` is a comma separated list of all the base URLs
* For running in Cloud Foundry, ensure a Config Server is bounded to the application. `VCAP_SERVICES` will be provided as an environment variables with the credentials to access the Config Server
  * If not running in Cloud Foundry but still want to connect to a Config Server via OAuth2, manually set the `VCAP_SERVICES` -- example value in `client/oauth2_test.go`

```go
package main

import (
	"fmt"
	"github.com/Piszmog/cloudconfigclient/client"
)

type File struct {
	Example Example `json:"example"`
}

type Example struct {
	Field string `json:"field"`
}

func main() {
	// To create a Client for a locally running Spring Config Server
	configClient, err := client.CreateLocalClient()
	// or to create a Client for a Spring Config Server in Cloud Foundry
	configClient, err := client.CreateCloudClient()
	if err != nil {
		panic(err)
	}
	var file File
	// Retrieves a 'temp1.json' from the Config Server's default branch in directory 'temp' and deserialize to File
	err = configClient.GetFile("temp", "temp1.json", &file)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", file)
	
	// Retrieves the configurations from the Config Server based on the application name and active profiles
	config, err := configClient.GetConfiguration("testApp", []string{"dev"})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", config)
}
```

## Config Client Creation
There are two type of clients that can be created. A local client for a locally running Config Server without security (OAuth2) 
and a cloud client for a Config Server running in a cloud environment.

### Local
To create a local client, call `client.CreateLocalClient()`. The client is configured with timeouts set and to use a pool of connections.

### Cloud
To create a cloud client, call `client.CreateCloudClient()`. The client is an OAuth2 client. The OAuth2 configurations are determined from the `VCAP_SERVICES` environment variable.

## Configurations
The Config Server allows the ability to retrieve configurations for an application. Only files that follow a strict naming 
convention will be loaded,

| File Name | 
| :---: |
|`application.{yml/properties}`|
|`application-{profile}.{yml/properties}`|
|`{application name}.{yml/properties}`|
|`{application name}-{profile}.{yml/properties}`|

The loaded configurations are in the following JSON format,

```json
{
  "name":"<name of application>",
  "profiles":"<profiles passed in request>",
  "label":"<GIT branch configurations loaded from>",
  "version":"<version>",
  "state":"<state>",
  "propertySources":[
    {
      "<propertySourceName>":{
        "name":"<property source name>",
        "source" : {
          "<source path in .properties format>":"<value>"
        }
      }
    }
  ]
}
```

To use the library to retrieve configurations, create a `client/ConfigClient` and 
invoke the method `GetConfiguration(applicationName string, profiles []string)`. The return will be the struct representation 
of the configuration JSON - `model.Configuration`.

## Resources
Spring's Config Server allows two ways to retrieve files from a backing repository.

| URL Path | 
| :---: |
|`/<appName>/<profiles>/<directory>/<file>?useDefaultLabel=true`|
|`/<appName>/<profiles>/<branch>/<directory>/<file>?useDefaultLabel=true`|

* When retrieving a file from the Config Server's default branch, the file must not exist at the root of the repository.

The functions available to retrieve resource files are, `GetFile(directory string, file string, interfaceType interface{}) error` and 
`GetFileFromBranch(branch string, directory string, file string, interfaceType interface{})`.

* The `interfaceTypee` is the object to deserialize the file to