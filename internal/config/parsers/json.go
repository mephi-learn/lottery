package parsers

import (
	"encoding/json"
)

var _ Parsers[struct{}] = (*JsonParser[struct{}])(nil)

type JsonParser[T any] struct{}

// Parse производит разбор файла.
func (p *JsonParser[T]) Parse(config string, result *T) error {
	return json.Unmarshal([]byte(config), &result)
}
