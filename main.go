package main

import (
	"fmt"
	"github.com/Piszmog/cloudconfigclient/client"
)

func main() {
	configClient, err := client.CreateLocalClient()
	if err != nil {
		panic(err)
	}
	// Retrieves the configurations from the Config Server based on the application name and active profiles
	config, err := configClient.GetConfiguration("testApp", []string{"dev"})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", config)
}
