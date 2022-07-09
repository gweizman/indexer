package content_parser

import (
	"fmt"
	"internal/format_parser"
)

const PrintSize int = 20

func printParse(file format_parser.ParsedFile) error {
	if len(file.Parsed) < PrintSize {
		fmt.Println(file.Path, string(file.Parsed)[:len(file.Parsed)])
	} else {
		fmt.Println(file.Path, string(file.Parsed)[:PrintSize])
	}

	return nil
}
