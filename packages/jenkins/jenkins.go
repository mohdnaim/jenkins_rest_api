package jenkins

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
)

// Details ...
type Details struct {
	URL      string
	Username string
	APIToken string
}

// JenkinsDetails ...
var JenkinsDetails = Details{}

// IsJobExist ...
func IsJobExist(projectName string) bool {
	fullURL := fmt.Sprintf("%sjob/%s", JenkinsDetails.URL, projectName)
	request, err := http.NewRequest("GET", fullURL, nil)
	request.SetBasicAuth(JenkinsDetails.Username, JenkinsDetails.APIToken)
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Println("project doesn't exist:", projectName)
		return false
	}
	if resp.StatusCode == 200 {
		// log.Println("project exists:", projectName)
		return true
	}
	log.Println("project doesn't exist:", projectName)
	return false
}

// CopyJenkinsJob ...
func CopyJenkinsJob(srcJob string, dstJob string) error {
	fullURL := fmt.Sprintf("%s%s?name=%s&mode=copy&from=%s", JenkinsDetails.URL, "createItem", dstJob, srcJob)
	request, err := http.NewRequest("POST", fullURL, nil)
	request.SetBasicAuth(JenkinsDetails.Username, JenkinsDetails.APIToken)
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	if resp.StatusCode == 200 {
		return nil
	}
	return fmt.Errorf("error copying job: %s", srcJob)
}

// DownloadConfigXML ...
func DownloadConfigXML(projectName string, dstFilename string) error {
	fullURL := fmt.Sprintf("%sjob/%s/config.xml", JenkinsDetails.URL, projectName)
	if DownloadFile(fullURL, dstFilename) == nil {
		return nil
	}
	return fmt.Errorf("error downloading file %s", fullURL)
}

// PostConfigXML ...
func PostConfigXML(projectName string, filename string) error {
	// read content of file
	data, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
		return err
	}

	fullURL := fmt.Sprintf("%sjob/%s/config.xml", JenkinsDetails.URL, projectName)
	request, err := http.NewRequest("POST", fullURL, data)
	request.SetBasicAuth(JenkinsDetails.Username, JenkinsDetails.APIToken)
	client := &http.Client{}

	// perform the request
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	if resp.StatusCode == 200 {
		// log.Printf("success postConfigXML: %s %s %s", projectName, filename, resp.Status)
		return nil
	}
	return fmt.Errorf("error postConfigXML: %s %s %s", projectName, filename, resp.Status)
}

// DownloadFile ...
func DownloadFile(url string, filepath string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	request, err := http.NewRequest("GET", url, nil)
	request.SetBasicAuth(JenkinsDetails.Username, JenkinsDetails.APIToken)
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// DownloadFileToBytes ...
func DownloadFileToBytes(url string) ([]byte, error) {
	request, err := http.NewRequest("GET", url, nil)
	request.SetBasicAuth(JenkinsDetails.Username, JenkinsDetails.APIToken)
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err) // writes to standard error
	}
	return bodyBytes, nil
}

// GetAllProjectNames ...
func GetAllProjectNames() []string {
	allProjectNames := make([]string, 0)
	allProjectsURL := fmt.Sprintf("%sapi/json?pretty=true", JenkinsDetails.URL)
	if respBytes, err := DownloadFileToBytes(allProjectsURL); err == nil {
		// json --> map
		var result map[string]interface{}
		json.Unmarshal([]byte(respBytes), &result)

		// if map has slice 'jobs', iterate over it
		if jobs, keyIsPresent := result["jobs"]; keyIsPresent && reflect.TypeOf(jobs).Kind() == reflect.Slice {
			jobs2 := result["jobs"].([]interface{})
			for _, job := range jobs2 {
				// log.Println(job)
				_map, _ := job.(map[string]interface{}) // assert to map
				allProjectNames = append(allProjectNames, _map["name"].(string))
			}
		}
	}
	return allProjectNames
}
