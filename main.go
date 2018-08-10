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
	localClient := client.CreateLocalClient()
	resourceClient := resource.CreateResourceClient(&localClient)
	var file File
	err := resourceClient.GetFile("temp", "temp1.json", &file)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", file)

	configurationClient := configuration.CreateconfigurationClient(&localClient)
	config, err := configurationClient.GetConfiguration("application", []string{"dev"})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", config)
}
