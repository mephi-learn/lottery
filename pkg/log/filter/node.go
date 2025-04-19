package filter

// filter предоставляет фильтр префиксного дерева для структурированных ключей.

import (
	"strings"
)

type Value = any

// node представляет дерево, где каждый узел является префиксом для своих дочерних узлов.
type node[V Value] struct {
	set      bool // Если это узел без значений, то false.
	value    V
	children map[string]node[V]
}

// cutoff описывает функцию итерации по набору значений.
type cutoff[T any] func(path T) (head string, tail T)

// get возвращает значение из [node]. Если либо полный путь или его родитель не был найден, будет возвращен false.
func get[T any, V Value](tree node[V], path T, iter cutoff[T]) (V, bool) {
	head, path := iter(path)

	child, ok := tree.children[head]
	if !ok {
		var val V
		return val, false
	}

	if child.children == nil {
		return child.value, child.set
	}

	return get[T, V](child, path, iter)
}

func (tree *node[V]) append(path []string, value V) {
	head, path := path[0], path[1:]

	child := tree.children[head]

	if len(path) > 0 {
		child.append(path, value)
	} else {
		child.value = value
		child.set = true
	}

	if tree.children == nil {
		tree.children = make(map[string]node[V])
	}

	tree.children[head] = child
}

// splitPath использует sep, чтобы получить первую и последнюю часть пути.
func splitPath(key string, sep rune) (string, string) {
	i := strings.IndexRune(key, sep)
	if i < 0 {
		return key, ""
	}

	return key[:i], key[i+1:]
}
