package main

import (
	"archive/zip"
	"fmt"
	"io"
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

	// extract zip file contents
	zf, err := zip.OpenReader(outFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer zf.Close()

	// Iterate through the files in the archive,
	for _, fl := range zf.File {
		flPath := filepath.Join(homeDir, instLoc, "java", fl.Name)
		fmt.Println("extracting file ", flPath)

		if fl.FileInfo().IsDir() {
			fmt.Println("creating directory")
			os.MkdirAll(flPath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(flPath), os.ModePerm); err != nil {
			panic(err)
		}

		dstFile, err := os.OpenFile(flPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fl.Mode())
		if err != nil {
			panic(err)
		}

		fileInArchive, err := fl.Open()
		if err != nil {
			panic(err)
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			panic(err)
		}
		fmt.Println("extraction completed")

		dstFile.Close()
		fileInArchive.Close()
	}
}
