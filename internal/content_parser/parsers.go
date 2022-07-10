package content_parser

import "internal/persistent_storage"

var content_parsers = [...]func(persistent_storage.ParsedFile) error{ctagsParse}
