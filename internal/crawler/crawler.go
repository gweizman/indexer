package crawler

import (
	"log"
	"os"
	"path/filepath"
)

func Crawl(path string, files chan string) {
	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			files <- path
			return nil
		})

	close(files)

	if err != nil {
		log.Println(err)
	}
}
