package crawler

import (
	"internal/persistent_storage"
	"log"
	"os"
	"path/filepath"
)

func Crawl(db *persistent_storage.IndexStorage, path string, files chan string) {
	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				tempErr := db.CreateDir(path)
				if tempErr != nil {
					log.Panic(tempErr)
				}
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
