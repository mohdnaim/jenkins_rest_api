package main

import (
	"fmt"
	"log"
	"strings"

	helpers "./packages/helpers"
	jenkins "./packages/jenkins"
)

func main() {
	// compulsory to set
	jenkinsURL := "http://127.0.0.1/"
	jenkinsUsername := "put your username here"
	jenkinsAPIToken := "put your API token here"

	jenkins.JenkinsDetails = jenkins.Details{jenkinsURL, jenkinsUsername, jenkinsAPIToken}
	xmlFolder := "xml"

	// 1. get all existing projects / jenkins jobs
	allProjectNames := jenkins.GetAllProjectNames()

	// 2. filter out projects that we want
	filteredProjectNames := make([]string, 0)
	for _, projectName := range allProjectNames {
		// do something

		// append to another slice based on condition
		if strings.HasPrefix(projectName, "prefix") {
			filteredProjectNames = append(filteredProjectNames, projectName)
		}
	}

	// 3. for each project, get its config.xml
	for _, projectName := range filteredProjectNames {
		xmlPath := fmt.Sprintf("%s/%s.xml", xmlFolder, projectName)
		if err := jenkins.DownloadConfigXML(projectName, xmlPath); err != nil {
			log.Println("error download config.xml for project:", projectName)
			continue // skip
		}
	}

	// 4. modify its config.xml
	files := helpers.GetFilenamesRecursively(xmlFolder)
	for _, xmlFile := range files {
		log.Println(xmlFile)
	}

	// 4b. rewrite config.xml

	// 5. http request POST updated config.xml
	for _, xmlFile := range files {
		tmpSlice := strings.Split(xmlFile, "/")
		projectName := tmpSlice[len(tmpSlice)-1]
		log.Println(projectName)

		if err := jenkins.PostConfigXML(projectName, xmlFile); err != nil {
			log.Println("error postconfigxml:", projectName)
		}
	}
}
