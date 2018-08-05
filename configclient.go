package main

import (
	"fmt"
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
	file := &File{}
	resourceClient := resource.CreateClient("http://localhost:8880")
	err := resourceClient.GetFileFromBranch("develop", "temp", "temp1.json", file)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", file)

	configClient := configuration.CreateClient("http://localhost:8880")
	configurations, err := configClient.GetConfiguration("exampleapp", []string{"dev"})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", configurations)
}
