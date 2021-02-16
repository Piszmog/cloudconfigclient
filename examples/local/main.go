package main

import (
	"fmt"
	"github.com/Piszmog/cloudconfigclient"
	"log"
	"net/http"
	"strings"
)

func main() {
	// ensure you have the Config Server running locally...

	client, err := cloudconfigclient.NewLocalClient(&http.Client{}, []string{"http://localhost:8888"})
	if err != nil {
		log.Fatalln(err)
	}

	// load a config file
	configuration, err := client.GetConfiguration("test-app", []string{"local"})
	if err != nil {
		log.Fatalln(err)
	}
	localProp, err := configuration.GetPropertySource("application-local.yml")
	if err != nil {
		log.Fatalln("failed to find local property file")
	}
	fmt.Printf("local property file: %+v\n", localProp)
	// handle all config properties
	configuration.HandlePropertySources(func(propertySource cloudconfigclient.PropertySource) {
		if strings.HasSuffix(propertySource.Name, "test-app.properties") {
			// TODO save off values
		}
	})

	// load a specific file (e.g. json/txt)
	var f map[string]string
	// if 'fooDir' has been added to 'searchPaths' in SCS v3.x, then pass "" (blank) for directory
	if err = client.GetFile("fooDir", "bar.json", &f); err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("file from default branch: %+v\n", f)

	// load a specific file (e.g. json/txt)
	var b map[string]string
	// if 'fooDir' has been added to 'searchPaths' in SCS v3.x, then pass "" (blank) for directory
	if err = client.GetFileFromBranch("develop", "fooDir", "bar.json", &b); err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("file from specific branch: %+v\n", b)
}
