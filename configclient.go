package main

import (
	"fmt"
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
	client := resource.CreateClient("http://localhost:8880")
	err := client.GetFileFromBranch("develop", "temp", "temp1.json", file)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", file)

	//configClient := configuration.Client{
	//	BaseUrl:    "http://localhost:8880",
	//	HttpClient: net.CreateDefaultHttpClient(),
	//}
	//configurations, err := configClient.GetConfiguration("exampleapp", []string{"dev"})
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("%+v", configurations)
}
