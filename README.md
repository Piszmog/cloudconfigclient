# Go Config Server Client
[![Build Status](https://travis-ci.org/Piszmog/cloudconfigclient.svg?branch=develop)](https://travis-ci.org/Piszmog/cloudconfigclient)
[![Go Report Card](https://goreportcard.com/badge/github.com/Piszmog/cloudconfigclient)](https://goreportcard.com/report/github.com/Piszmog/cloudconfigclient)
[![GitHub release](https://img.shields.io/github/release/Piszmog/cloudconfigclient.svg)](https://github.com/Piszmog/cloudconfigclient/releases/latest)

Go library for Spring Config Server. Inspired by the Java library [Cloud Config Client](https://github.com/Piszmog/cloud-config-client)

## Description
Spring's Config Server provides way to externalize configurations of applications. Spring's
[Spring Cloud Config Client](https://github.com/spring-cloud/spring-cloud-config/tree/master/spring-cloud-config-client)
can be used to load the base configurations an application requires to function.

This library provides clients the ability to load Configurations and Files from the Config Server

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

To use the library to retrieve configurations, use a client of `configuration.Configuration` called `configuration.Client` and 
invoke the method `GetConfiguration(applicationName string, profiles []string)`. The return will be the struct representation 
of the configuration JSON - `model.Configuration`.

## Resources
