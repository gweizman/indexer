package content_parser

import (
	"internal/persistent_storage"
	"log"
	"regexp"

	ctags "github.com/sourcegraph/go-ctags"
)

func ctagsParse(file persistent_storage.ParsedFile) error {
	p, err := ctags.New(ctags.Options{
		Bin: "ctags",
	})

	if err != nil {
		return err
	}

	re := regexp.MustCompile(`\r?\n`)
	better_content := re.ReplaceAllString(string(file.Parsed), "\n")

	got, err := p.Parse(file.Path, []byte(better_content))
	if err != nil {
		return err
	}
	for _, g := range got {
		log.Println(g)
	}

	return nil
}
