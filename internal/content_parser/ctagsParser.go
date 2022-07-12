package content_parser

import (
	"internal/persistent_storage"
	"regexp"

	ctags "github.com/sourcegraph/go-ctags"
)

func ctagsParse(file persistent_storage.StoredFile, content []byte) error {
	p, err := ctags.New(ctags.Options{
		Bin: "ctags",
	})

	if err != nil {
		return err
	}

	re := regexp.MustCompile(`\r?\n`)
	better_content := re.ReplaceAllString(string(content), "\n")

	got, err := p.Parse(file.Name, []byte(better_content))
	if err != nil {
		return err
	}
	for _, g := range got {
		file.InsertDefinition(g.Name, g.Line, g.Kind, g.Language, g.Parent, g.ParentKind, g.Pattern, g.Signature, g.FileLimited)
	}

	return nil
}
