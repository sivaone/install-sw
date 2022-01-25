package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/codeclysm/extract/v3"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	fileUrl := "https://www.learningcontainer.com/wp-content/uploads/2020/05/sample-zip-file.zip"
	instLoc := "temp"

	urlStruct, err := url.Parse(fileUrl)
	if err != nil {
		log.Fatalf("unable to parse url: %v", err)
	}

	homeDir, err := os.UserHomeDir()

	path := urlStruct.Path
	segments := strings.Split(path, "/")
	fileName := segments[len(segments)-1]
	outFileName := filepath.Join(homeDir, instLoc, fileName)

	// create blank file with given filename
	file, err := os.Create(outFileName)
	if err != nil {
		log.Fatalf("unable to create file on local disk: %v", err)
	}

	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			req.URL.Opaque = req.URL.Path
			return nil
		},
	}

	// write content from url to file
	resp, err := client.Get(fileUrl)
	if err != nil {
		log.Fatalf("unable to retrieve file from url: %v", err)
	}
	defer resp.Body.Close()

	size, err := io.Copy(file, resp.Body)
	if err != nil {
		log.Fatalf("unable to write response body to file: %v", err)
	}
	defer file.Close()

	fmt.Printf("downloaded file %s with size %d \n", outFileName, size)

	readFile, err := ioutil.ReadFile(outFileName)
	if err != nil {
		panic(err)
	}

	flPath := filepath.Join(homeDir, instLoc, "java")
	buffer := bytes.NewBuffer(readFile)
	extract.Archive(context.Background(), buffer, flPath, nil)

	fmt.Println("extraction is completed")
}
