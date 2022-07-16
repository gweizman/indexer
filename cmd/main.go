package main

import (
	"flag"
	"internal/content_parser"
	"internal/crawler"
	"internal/format_parser"
	"internal/persistent_storage"
	"log"
	"time"
)

func addProject(dir string, name string, db persistent_storage.Db) {
	storage, err := persistent_storage.NewStorage(name, dir, db)
	if err != nil {
		log.Panicf("Failed creating storage: %s", err.Error())
	}

	files_names := make(chan string, 3000)
	files := make(chan persistent_storage.StoredFile, 3000)
	done := make(chan int)

	go crawler.Crawl(storage, dir, files_names)
	go format_parser.Parse(storage, files_names, files, 100)
	go content_parser.Parse(storage, files, 100, done)

	<-done
}

func main() {
	start := time.Now()
	initSchema := flag.Bool("initSchema", false, "Init the schemas")
	indexProject := flag.Bool("indexProject", false, "Index the project") // TODO: Also add a flag for project path/etc.
	flag.Parse()

	db, err := persistent_storage.CreateDb()
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	if *initSchema {
		err = db.InitSchemas()
		if err != nil {
			log.Panic(err)
		}
	}
	if *indexProject {
		addProject("D:\\Work\\hellogitworld", "test_project", *db)
	}

	elapsed := time.Since(start)
	log.Printf("Indexing done. Took %s.", elapsed)
}
