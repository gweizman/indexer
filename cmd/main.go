package main

import (
	"internal/content_parser"
	"internal/crawler"
	"internal/format_parser"
	"log"
	"time"
)

func addProject(dir string, name string) {
	files := make(chan string, 3000)
	contents := make(chan format_parser.ParsedFile, 3000)
	done := make(chan int)

	go crawler.Crawl(dir, files)
	go format_parser.Parse(files, contents, 100)
	go content_parser.Parse(contents, 100, done)

	<-done
}

func main() {
	start := time.Now()

	addProject("D:\\Work\\hellogitworld", "hello")

	elapsed := time.Since(start)
	log.Printf("took %s", elapsed)
}
