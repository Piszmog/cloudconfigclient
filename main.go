package main

import (
	"fmt"
	"github.com/Piszmog/cloudconfigclient/client"
	"github.com/Piszmog/cloudconfigclient/configuration"
	"github.com/Piszmog/cloudconfigclient/resource"
)

type File struct {
	Example Example `json:"example"`
}

type Example struct {
	Field3 string `json:"field3"`
}

func main() {
	//configClient := client.CreateLocalClient()
	configClient, err := client.CreateCloudClient()
	if err != nil {
		panic(err)
	}
	resourceClient := resource.CreateResourceClient(configClient)
	var file File
	err = resourceClient.GetFile("temp", "temp1.json", &file)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", file)

	configurationClient := configuration.CreateConfigurationClient(configClient)
	config, err := configurationClient.GetConfiguration("application", []string{"dev"})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", config)
}
