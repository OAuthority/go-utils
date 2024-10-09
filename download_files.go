// This script assists with downloading all of the images from
// a MediaWiki wiki such as Fandom (who will not provide image dumps)
// Scraping Fandom is against their ToU, but who cares. Use this wisely.
//
// This is free software licensed under the MIT License
// See LICENSE
// (c) OAuthority 2024

package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// general structure of the api response that will be returned
// inside each of the query is a reference to the individual structure ImageInfo
type ApiResponse struct {
	Batchcomplete bool `json:"batchcomplete"`
	Query         struct {
		Allimages []ImageInfo `json:"allimages"`
	} `json:"query"`
}

// struct to represent the structure of each individual image returned
// by the api
type ImageInfo struct {
	Name           string `json:"name"`
	Timestamp      string `json:"timestamp"`
	URL            string `json:"url"`
	DescriptionURL string `json:"descriptionurl"`
	NS             int    `json:"ns"`
	Title          string `json:"title"`
}

func main() {

	// the api endpoint we will call to get a list of all of the files
	// on the wiki
	apiEndpoint := flag.String("api", "", "The API endpoint to get a list of images from.")

	// where should we write the files to? Doesn't matter if it exists or not, if not we will create it
	outFolder := flag.String("output", "", "The folder to write the files to")

	flag.Usage = func() {

		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", "download_files.go")

		fmt.Println("This script aids in downloading all of the images from a MediaWiki wiki to a local folder on your hard drive. ")

		flag.PrintDefaults()
	}

	flag.Parse()

	if *apiEndpoint == "" || *outFolder == "" {
		fmt.Println("You must provide both the -api flag and the -output flag")
		flag.Usage()
		os.Exit(1)
	}

	allImages, err := getAllImages(*apiEndpoint)

	if err != nil {
		log.Fatalf("Error getting images: %v", err)
	}

	// check if the directory exists and if not, create it
	// this is a bit backward, first we try to stat the directory and
	// if we receive an error then we know it doesn't exist.
	if _, err := os.Stat(*outFolder); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(*outFolder, os.ModePerm)
		if err != nil {
			log.Fatalf("Error creating directory: %v", err)
		}
	}

	for _, img := range allImages {
		downloadImage(*outFolder, img)
	}
}

func getAllImages(apiEndpoint string) ([]ImageInfo, error) {

	if !strings.HasSuffix(apiEndpoint, "api.php") {
		apiEndpoint = strings.TrimSuffix(apiEndpoint, "/") + "/api.php"
	}

	// construct the base url for the api request using the url that the user gave
	apiBaseUrl := fmt.Sprintf("%s?action=query&format=json&list=allimages&formatversion=2", apiEndpoint)

	params := url.Values{}
	params.Add("ailimit", "500")

	fullUrl := apiBaseUrl + "&" + params.Encode()

	resp, err := http.Get(fullUrl)

	if err != nil {
		return nil, fmt.Errorf("error making GET request: %v", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var result ApiResponse

	err = json.Unmarshal(body, &result)

	if err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}

	// Return the original JSON
	return result.Query.Allimages, nil
}

// download an image to the specified directory
// doesn't do much other than that. The check for the directories existence
// is done a bit earlier in the main() function
func downloadImage(dir string, img ImageInfo) {
	url := img.URL
	var imageName string = img.Name

	response, err := http.Get(url)

	if err != nil {
		log.Fatalf("Error fetching file: %v", err)
	}

	defer response.Body.Close()

	file, err := os.Create(dir + "/" + imageName)

	if err != nil {
		log.Fatalf("Error opening file for write: %v", err)
	}

	_, err = io.Copy(file, response.Body)

	if err != nil {
		log.Fatalf("Error writing file to disk: %v", err)
	}

	fmt.Printf("Successfully written %s to disk\n", imageName)
}
