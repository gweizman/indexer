module newgrok

go 1.18

require internal/crawler v0.0.0-00010101000000-000000000000

replace internal/crawler => ./internal/crawler

require internal/format_parser v0.0.0-00010101000000-000000000000

replace internal/format_parser => ./internal/format_parser

require internal/content_parser v0.0.0-00010101000000-000000000000

replace internal/content_parser => ./internal/content_parser

require internal/persistent_storage v0.0.0-00010101000000-000000000000

replace internal/persistent_storage => ./internal/persistent_storage

require (
	github.com/go-chi/chi/v5 v5.0.7 // indirect
	github.com/go-chi/render v1.0.1 // indirect
	github.com/gocql/gocql v1.2.0 // indirect
	github.com/golang/snappy v0.0.3 // indirect
	github.com/hailocab/go-hostpool v0.0.0-20160125115350-e80d13ce29ed // indirect
	github.com/neo4j/neo4j-go-driver/v4 v4.4.3 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sourcegraph/go-ctags v0.0.0-20220611154803-db463692f037 // indirect
	github.com/stevenferrer/solr-go v0.3.2 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
)
