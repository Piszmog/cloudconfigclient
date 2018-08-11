package main

import (
	"fmt"
	"github.com/Piszmog/cloudconfigclient/client"
)

type File struct {
	Example Example `json:"example"`
}

type Example struct {
	Field3 string `json:"field3"`
}

func main() {
	//configClient, err := client.CreateLocalClient()
	configClient, err := client.CreateCloudClient()
	if err != nil {
		panic(err)
	}
	var file File
	err = configClient.GetFile("temp", "temp1.json", &file)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", file)

	config, err := configClient.GetConfiguration("testapp", []string{"dev"})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", config)
}
