# Go Config Server Client

[![Go Reference](https://pkg.go.dev/badge/github.com/Piszmog/cloudconfigclient.svg)](https://pkg.go.dev/github.com/Piszmog/cloudconfigclient)
[![Build Status](https://github.com/Piszmog/cloudconfigclient/workflows/Go/badge.svg)](https://github.com/Piszmog/cloudconfigclient/workflows/Go/badge.svg)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=Piszmog_cloudconfigclient&metric=alert_status)](https://sonarcloud.io/dashboard?id=Piszmog_cloudconfigclient)
[![Coverage Status](https://coveralls.io/repos/github/Piszmog/cloudconfigclient/badge.svg?branch=main)](https://coveralls.io/github/Piszmog/cloudconfigclient?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/Piszmog/cloudconfigclient)](https://goreportcard.com/report/github.com/Piszmog/cloudconfigclient)
[![GitHub release](https://img.shields.io/github/release/Piszmog/cloudconfigclient.svg)](https://github.com/Piszmog/cloudconfigclient/releases/latest)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Go library for Spring Config Server. Inspired by the Java
library [Cloud Config Client](https://github.com/Piszmog/cloud-config-client).

`go get github.com/Piszmog/cloudconfigclient/v2`

#### V2 Migration

See [V2 Migration](https://github.com/Piszmog/cloudconfigclient/wiki/V2-Migration) for details on how to migrate from V1
to V2

## Description

Spring's Config Server provides way to externalize configurations of applications. Spring's
[Spring Cloud Config Client](https://github.com/spring-cloud/spring-cloud-config/tree/master/spring-cloud-config-client)
can be used to load the base configurations that an application requires to properly function.

This library provides clients the ability to load Configurations and Files from the Config Server.

### Compatibility

This library is compatible with versions of Spring Config Server greater than or equal to `1.4.x.RELEASE`. Prior
versions of the Config Server do not provide the endpoint necessary to retrieve files for the Config Server's default
branch.

#### Spring Cloud Config Server v3.x

Since Spring Cloud Services v3.0, the service name in `VCAP_SERVICES` has changed from `p-config-server` to
be `p.config-server`.

To help mitigate migration difficulties, `cloudconfigclient.New(cloudconfigclient.DefaultCFService())` will first search
for the service `p.config-server` (v3.x). If the v3.x service could not be found,
`p-config-server` (v2.x) will be search for.

See [Spring Cloud Services Differences](https://docs.pivotal.io/spring-cloud-services/3-1/common/config-server/managing-service-instances.html#differences-between-3-0-and-earlier)
for more details.

## Example Usage

Below is an example usage of the library to retrieve a file from the Config Server and to retrieve the application's
configurations

* For local config client, there are two options (`Option`) the create a client
    1. Call `LocalEnv()`. Set the environment variable `CONFIG_SERVER_URLS`. It is a comma separated list of all the
       base URLs
    2. Call `Local(baseUrls ...string)`. Provide the array of base URLs of Config Servers.
* For running in Cloud Foundry, ensure a Config Server is bounded to the application. `VCAP_SERVICES` will be provided
  as an environment variables with the credentials to access the Config Server
* For connecting to a Config Server via OAuth2 and not deployed to Cloud Foundry, an OAuth2 Client can be created
  with `OAuth2(baseURL string, clientId string, secret string, tokenURI string)`

```go
package main

import (
	"fmt"
	"github.com/Piszmog/cloudconfigclient/v2"
	"net/http"
)

type File struct {
	Example Example `json:"example"`
}

type Example struct {
	Field string `json:"field"`
}

func main() {
	// To create a Client for a locally running Spring Config Server
	configClient, err := cloudconfigclient.New(cloudconfigclient.LocalEnv(&http.Client{}))
	// Or
	configClient, err = cloudconfigclient.New(cloudconfigclient.Local(&http.Client{}, "http://localhost:8888"))
	// or to create a Client for a Spring Config Server in Cloud Foundry
	configClient, err = cloudconfigclient.New(cloudconfigclient.DefaultCFService())
	// or to create a Client for a Spring Config Server with OAuth2
	configClient, err = cloudconfigclient.New(cloudconfigclient.OAuth2("config server uri", "client id", "client secret",
		"access token uri"))
	// or a combination of local, Cloud Foundry, and OAuth2
	configClient, err = cloudconfigclient.New(
		cloudconfigclient.Local(&http.Client{}, "http://localhost:8888"),
		cloudconfigclient.DefaultCFService(),
		cloudconfigclient.OAuth2("config server uri", "client id", "client secret", "access token uri"),
	)

	if err != nil {
		fmt.Println(err)
		return
	}
	var file File
	// Retrieves a 'temp1.json' from the Config Server's default branch in directory 'temp' and deserialize to File
	err = configClient.GetFile("temp", "temp1.json", &file)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", file)
	// Retrieves a 'temp2.txt' from the Config Server's default branch in directory 'temp' as a byte slice ([]byte)
	b, err := configClient.GetFileRaw("temp", "temp2.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))

	// Retrieves the configurations from the Config Server based on the application name and active profiles
	config, err := configClient.GetConfiguration("testApp", "dev")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v", config)
}
```

#### VCAP_SERVICES

When an application is deployed to Cloud Foundry, services can be bounded to the application. When a service is bounded
to an application, the application will have the necessary connection information provided in the environment
variable `VCAP_SERVICES`.

Structure of the `VCAP_SERVICES` value

```json
{
  "<service type :: e.g. p-config-server>": [
    {
      "name": "<the service name>",
      "instance_name": "<service name>",
      "binding_name": "<bounded name of the service>",
      "credentials": {
        "uri": "<URI of the service :: used to connect to the service>",
        "client_secret": "<OAuth2 client secret>",
        "client_id": "<OAuth2 client id>",
        "access_token_uri": "<OAuth2 token URI>"
      },
      ...
    }
  ]
}
```

##### CredHub Reference

Newer versions of PCF (>=2.6) may have services that use a CredHub Reference to store credential information.

When viewing the Environment Variables of an application via the UI, the credentials may appear as the following

```json
{
  "credentials": {
    "credhub-ref": "/c/example-service-broker/example-service/faa677f5-25cd-4f1e-8921-14a9d5ab48b8/credentials"
  }
}
```

When the application starts up, the `credhub-ref` is replaced with the actual credential values that application will
need to connect to the service.

## Configurations

The Config Server allows the ability to retrieve configurations for an application. Only files that follow a strict
naming convention will be loaded,

| File Name | 
| :---: |
|`application.{yml/properties}`|
|`application-{profile}.{yml/properties}`|
|`{application name}.{yml/properties}`|
|`{application name}-{profile}.{yml/properties}`|

The loaded configurations are in the following JSON format,

```json
{
  "name": "<name of application>",
  "profiles": "<profiles passed in request>",
  "label": "<GIT branch configurations loaded from>",
  "version": "<version>",
  "state": "<state>",
  "propertySources": [
    {
      "<propertySourceName>": {
        "name": "<property source name>",
        "source": {
          "<source path in .properties format>": "<value>"
        }
      }
    }
  ]
}
```

To use the library to retrieve configurations, create a `Client` and invoke the
method `GetConfiguration(applicationName string, profiles ...string)`. The return will be the struct representation of
the configuration JSON - `client.Configuration`.

## Resources

Spring's Config Server allows two ways to retrieve files from a backing repository.

| URL Path | 
| :---: |
|`/<appName>/<profiles>/<directory>/<file>?useDefaultLabel=true`|
|`/<appName>/<profiles>/<branch>/<directory>/<file>`|

* When retrieving a file from the Config Server's default branch, the file must not exist at the root of the repository.
* If the `directory` is in the `searchPath`, it does not have to be specified (depending on SCCS version)

The functions available to retrieve resource files
are `GetFile(directory string, file string, interfaceType interface{})` and
`GetFileFromBranch(branch string, directory string, file string, interfaceType interface{})`. To retrieve the data from
the files, the functions available are `GetFileRaw(directory string, file string)` and
`GetFileFromBranchRaw(branch string, directory string, file string)`

* The `interfaceType` is the object to deserialize the file to

### Spring Cloud Config Server v3.x Changes

The following is only for certain versions of SCCS v3.x. If a file is not being found by the client, the following may
be true.

SCCS v3.x slightly changed how files are retrieved. If the Config Server specified a directory in the `searchPaths`, the
path should be excluded from the `GetFile(..)` invocation.

For example if `common` has been specified in the `searchPaths` and the file `common/foo.txt` needs to be retrieved,
then the `directory` to provide to `GetFile(..)`
should be `""` (blank).

This differs with SCS v2.x where the directory in `searchPaths` did not impact the `directory` provided
to `GetFile(..)` (e.g. to retrieve file `common/foo.txt`,
`directory` would be `"common"`).
