package content_parser

import (
	"internal/persistent_storage"
)

const PrintSize int = 20

func printParse(file *persistent_storage.StoredFile) error {
	// if len(file.Parsed) < PrintSize {
	// 	fmt.Println(file.Path, string(file.Parsed)[:len(file.Parsed)])
	// } else {
	// 	fmt.Println(file.Path, string(file.Parsed)[:PrintSize])
	// }

	return nil
}
