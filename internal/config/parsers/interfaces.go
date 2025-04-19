package parsers

type Parsers[T any] interface {
	Parse(config string, data *T) error
}
