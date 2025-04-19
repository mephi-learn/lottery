package parsers

import (
	"gopkg.in/yaml.v3"
)

var _ Parsers[struct{}] = (*YamlParser[struct{}])(nil)

type YamlParser[T any] struct{}

// Parse производит разбор файла.
func (p *YamlParser[T]) Parse(config string, result *T) error {
	return yaml.Unmarshal([]byte(config), result)
}
