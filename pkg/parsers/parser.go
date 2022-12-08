package parsers

type NewParserFunc func() (Parser, error)

type Parser interface {
	Parse() (any, error)
}
