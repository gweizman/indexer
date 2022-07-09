package content_parser

import "internal/format_parser"

var content_parsers = [...]func(format_parser.ParsedFile) error{ctagsParse}
