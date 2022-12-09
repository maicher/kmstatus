package parsers

type Parser interface {
	Parse() (any, error)
}
