package main

import (
	"internal/content_parser"
	"internal/crawler"
	"internal/format_parser"
	"internal/persistent_storage"
	"log"
	"time"
)

func addProject(dir string, name string) {
	storage, err := persistent_storage.NewStorage(name, dir)
	if err != nil {
		log.Panic("Failed creating storage")
	}

	files_names := make(chan string, 3000)
	files := make(chan persistent_storage.StoredFile, 3000)
	done := make(chan int)

	go crawler.Crawl(dir, files_names)
	go format_parser.Parse(storage, files_names, files, 100)
	go content_parser.Parse(storage, files, 100, done)

	<-done
}

func main() {
	start := time.Now()

	addProject("D:\\Work\\test", "hello")

	elapsed := time.Since(start)
	log.Printf("took %s", elapsed)
}
