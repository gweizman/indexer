package format_parser

import (
	"bytes"
	"internal/persistent_storage"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

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

func worker(storage *persistent_storage.IndexStorage, files chan string, results chan persistent_storage.StoredFile, done chan int) {
	for file_name := range files {
		func() {
			fileObject, err := storage.CreateFile(file_name)
			if err != nil {
				log.Panic("Failed creating file object")
			}

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

			fileObject.AddContent(content)
			fileObject.AddContentVersion(tika(filepath.Base(file_name), content), "Tika")

			results <- fileObject
		}()
	}

	done <- 1
}

func Parse(storage *persistent_storage.IndexStorage, files chan string, contents chan persistent_storage.StoredFile, workerCount int) {
	done := make(chan int, workerCount)

	for w := 0; w < workerCount; w++ {
		go worker(storage, files, contents, done)
	}

	for j := 0; j < workerCount; j++ {
		<-done
	}

	close(contents)
}
