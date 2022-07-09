package content_parser

import (
	"internal/format_parser"
	"log"
)

func worker(contents chan format_parser.ParsedFile, done chan int) {
	for i := range contents {
		for _, content_parser := range content_parsers {
			err := content_parser(i)
			if err != nil {
				log.Print(err)
			}
		}
	}

	done <- 1
}

func Parse(contents chan format_parser.ParsedFile, workerCount int, done chan int) {
	myDone := make(chan int, workerCount)

	for w := 0; w < workerCount; w++ {
		go worker(contents, myDone)
	}

	for j := 0; j < workerCount; j++ {
		<-myDone
	}

	done <- 1
}
