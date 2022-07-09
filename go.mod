module example/hello

go 1.18

require internal/crawler v0.0.0-00010101000000-000000000000

replace internal/crawler => ./internal/crawler

require internal/format_parser v0.0.0-00010101000000-000000000000

replace internal/format_parser => ./internal/format_parser

require internal/content_parser v0.0.0-00010101000000-000000000000

require github.com/sourcegraph/go-ctags v0.0.0-20220611154803-db463692f037 // indirect

replace internal/content_parser => ./internal/content_parser
