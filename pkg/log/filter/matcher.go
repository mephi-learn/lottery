package filter

import "strings"

type noval = struct{}

type StringMatcher struct {
	node[struct{}]

	sep rune
}

// NewStringMatcher создает фильтры для соответсвия ключей.
func NewStringMatcher(keys []string, sep rune) StringMatcher {
	var root node[noval]

	for _, key := range keys {
		path := strings.Split(key, string(sep))
		root.append(path, noval{})
	}

	return StringMatcher{node: root, sep: sep}
}

// Match возвращает true если путь найден в дереве.
// Если есть более широкий путь, тоже вернется true.
//
// Пример:
//
//	tree: [ a ]
//	 - Match(a) == true
//	 - Match(a, b) == true
//	 - Match(x) == false
//
//	tree: [ a.b ]
//	 - Match(a) == false
//	 - Match(a, b) == true
//	 - Match(a, b, x) == true
//	 - Match(a, c) == false
func (m StringMatcher) Match(path string) bool {
	_, ok := get(m.node, path, m.cutoff)
	return ok
}

func (m StringMatcher) cutoff(path string) (string, string) {
	return splitPath(path, m.sep)
}
