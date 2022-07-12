package content_parser

import (
	"internal/persistent_storage"
	"log"
)

func worker(contents chan persistent_storage.StoredFile, done chan int) {
	for file := range contents {
		for _, content := range file.Content {
			data, canBeCode := content.Data, content.CanBeCode
			if canBeCode {
				for _, content_parser := range content_parsers {
					err := content_parser(file, data)
					if err != nil {
						log.Print(err)
					}
				}
			}
		}
	}

	done <- 1
}

func Parse(storage *persistent_storage.IndexStorage, contents chan persistent_storage.StoredFile, workerCount int, done chan int) {
	myDone := make(chan int, workerCount)

	for w := 0; w < workerCount; w++ {
		go worker(contents, myDone)
	}

	for j := 0; j < workerCount; j++ {
		<-myDone
	}

	done <- 1
}
