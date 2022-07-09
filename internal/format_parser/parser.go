package format_parser

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type ParsedFile struct {
	Path   string
	Plain  []byte
	Parsed []byte
}

func tika(filename string, content []byte) []byte {
	req, err := http.NewRequest("PUT", "http://localhost:9998/tika/", bytes.NewBuffer(content))
	if err != nil {
		log.Println(err)
		return nil
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("filename", filename)
	req.Header.Set("Accept", "text/plain")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return nil
	}

	return body
}

func worker(files chan string, results chan ParsedFile, done chan int) {
	for file_name := range files {
		func() {
			file, err := os.Open(file_name)
			if err != nil {
				log.Println(err)
				return
			}
			defer file.Close()

			content, err := ioutil.ReadAll(file)
			if err != nil {
				log.Println(err)
				return
			}

			results <- ParsedFile{
				file_name,
				[]byte(content),
				tika(filepath.Base(file_name), content),
			}
		}()
	}

	done <- 1
}

func Parse(files chan string, contents chan ParsedFile, workerCount int) {
	done := make(chan int, workerCount)

	for w := 0; w < workerCount; w++ {
		go worker(files, contents, done)
	}

	for j := 0; j < workerCount; j++ {
		<-done
	}

	close(contents)
}
