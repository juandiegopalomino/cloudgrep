package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//generate the markdown files using the template files
func main() {

	//generate the data for the template
	data := make(map[string]any)

	//read supported resources
	suportedResourceFile, err := ioutil.ReadFile("./pkg/provider/aws/zz_integration_stats.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}
	var supportedResources []string
	err = json.Unmarshal(suportedResourceFile, &supportedResources)
	if err != nil {
		log.Fatal(err)
	}
	data["supportedResources"] = supportedResources

	//add the config yaml file
	configFile, err := ioutil.ReadFile("./pkg/config/config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	data["configFile"] = string(configFile)

	//read all the template files and generate the markdown files
	tmplFiles, err := filepath.Glob("./*.md.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	for _, tmplFile := range tmplFiles {
		// Create the output file
		targetFile := strings.Replace(tmplFile, ".md.tmpl", ".md", 1)
		f, err := os.Create(targetFile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		//generate the markdown file from the template
		log.Default().Printf("Writing file %v", targetFile)
		err = template.Must(template.ParseFiles(tmplFile)).Execute(f, data)
		if err != nil {
			log.Fatal(err)
		}
	}

}
