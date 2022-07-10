package content_parser

import (
	"fmt"
	"internal/persistent_storage"
)

const PrintSize int = 20

func printParse(file persistent_storage.ParsedFile) error {
	if len(file.Parsed) < PrintSize {
		fmt.Println(file.Path, string(file.Parsed)[:len(file.Parsed)])
	} else {
		fmt.Println(file.Path, string(file.Parsed)[:PrintSize])
	}

	return nil
}
