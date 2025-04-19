package filter

import (
	"fmt"
	"strings"
)

type Key = fmt.Stringer

// forward возвращает первый ключ.
func forward[K Key](keys []K) (string, []K) {
	if len(keys) > 0 {
		return keys[0].String(), keys[1:]
	}

	return "", nil
}

// backward возвращает последний ключ.
func backward[K Key](keys []K) (string, []K) {
	if last := len(keys) - 1; last >= 0 {
		return keys[last].String(), keys[:last]
	}

	return "", nil
}

type Filter[K Key, V Value] struct {
	node[V]

	cutoff cutoff[[]K]
}

// NewBackwardFilter создает фильтр для прохождения пути в порядке "root last".
func NewBackwardFilter[K Key, V Value](pairs map[string]V, sep rune) Filter[K, V] {
	var root node[V]

	for key, val := range pairs {
		path := strings.Split(key, string(sep))
		root.append(path, val)
	}

	return Filter[K, V]{node: root, cutoff: backward[K]}
}

// Get возвращает значение.
func (f Filter[K, V]) Get(path []K) (V, bool) {
	return get(f.node, path, f.cutoff)
}
