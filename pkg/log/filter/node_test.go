package filter

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

type (
	testKey = *bytes.Buffer
	testVal = *struct{}
)

func testPath(keys ...string) []*bytes.Buffer {
	path := make([]*bytes.Buffer, len(keys))

	for i, key := range keys {
		path[i] = bytes.NewBufferString(key)
	}

	return path
}

func TestGetEmptyTree(t *testing.T) {
	t.Parallel()

	var tree node[testVal]

	tests := []struct {
		name string
		path []testKey
	}{
		{"nil", nil},
		{"empty", []testKey{}},
		{"single", testPath("a")},
		{"nested", testPath("a", "b")},
	}

	// все тест кейсы будут возвращать (nil, false)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			val, ok := get(tree, nil, forward[testKey])
			require.False(t, ok)
			require.Nil(t, val)
		})
	}
}

func TestGetSingle(t *testing.T) {
	t.Parallel()

	expected := &struct{}{}

	var tree node[testVal]
	tree.append([]string{"a"}, expected)

	tests := []struct {
		name string
		path []testKey
		val  testVal
		ok   bool
	}{
		{
			name: "nil",
			path: nil,
			ok:   false,
		},
		{
			name: "empty",
			path: []testKey{},
			ok:   false,
		},
		{
			name: "single",
			path: testPath("a"),
			ok:   true,
		},
		{
			name: "nested",
			path: testPath("a", "b"),
			ok:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			val, ok := get(tree, test.path, forward[testKey])
			require.Equal(t, test.ok, ok)

			if ok {
				require.Equal(t, expected, val)
			} else {
				require.Nil(t, val)
			}
		})
	}
}

func TestGetNested(t *testing.T) {
	t.Parallel()

	expected := &struct{}{}

	var tree node[testVal]
	tree.append([]string{"a", "b"}, expected)

	tests := []struct {
		name string
		path []testKey
		val  testVal
		ok   bool
	}{
		{
			name: "nil",
			path: nil,
			ok:   false,
		},
		{
			name: "empty",
			path: []testKey{},
			ok:   false,
		},
		{
			name: "single",
			path: testPath("a"),
			ok:   false,
		},
		{
			name: "nested",
			path: testPath("a", "b"),
			ok:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			val, ok := get(tree, test.path, forward[testKey])
			require.Equal(t, test.ok, ok)

			if ok {
				require.Equal(t, expected, val)
			} else {
				require.Nil(t, val)
			}
		})
	}
}

func TestAppendSingle(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	var tree node[testVal]
	tree.append([]string{"a"}, nil)

	// корень существует и установлен.
	require.Contains(tree.children, "a")
	require.NotNil(tree.children["a"])

	a := tree.children["a"]
	require.Nil(a.children)
	require.Nil(a.value)
	require.True(a.set)
}

func TestAppendNeighbors(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	var tree node[testVal]
	tree.append([]string{"a", "b"}, nil)
	tree.append([]string{"a", "c"}, nil)

	require.Contains(tree.children, "a")
	require.NotNil(tree.children["a"])

	// корень существует, но не установлен.
	a := tree.children["a"]
	require.Nil(a.value)
	require.False(a.set)

	// ребенок b существует и установлен.
	require.Contains(a.children, "b")
	require.NotNil(a.children["b"])

	b := a.children["b"]
	require.Nil(b.children)
	require.Nil(b.value)
	require.True(b.set)

	// ребенок c существует и установлен.
	require.Contains(a.children, "c")
	require.NotNil(a.children["c"])

	c := a.children["c"]
	require.Nil(c.children)
	require.Nil(c.value)
	require.True(c.set)
}

func TestAppendNested(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	var tree node[testVal]
	tree.append([]string{"a"}, nil)
	tree.append([]string{"a", "b"}, nil)

	// корень существует и установлен.
	require.Contains(tree.children, "a")
	require.NotNil(tree.children["a"])

	a := tree.children["a"]
	require.Nil(a.value)
	require.True(a.set)

	// ребенок b существует и установлен.
	require.Contains(a.children, "b")
	require.NotNil(a.children["b"])

	b := a.children["b"]
	require.Nil(b.children)
	require.Nil(b.value)
	require.True(b.set)
}

func TestAppendNestedReverse(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	var tree node[testVal]
	tree.append([]string{"a", "b"}, nil)
	tree.append([]string{"a"}, nil)

	// корень существует и установлен.
	require.Contains(tree.children, "a")
	require.NotNil(tree.children["a"])

	a := tree.children["a"]
	require.Nil(a.value)
	require.True(a.set)

	// ребенок b существует и установлен.
	require.Contains(a.children, "b")
	require.NotNil(a.children["b"])

	b := a.children["b"]
	require.Nil(b.children)
	require.Nil(b.value)
	require.True(b.set)
}

func TestSplitPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		path string
		sep  rune
		exp  []string
	}{
		{"", '.', []string{"", ""}},
		{".", '.', []string{"", ""}},
		{".x", '.', []string{"", "x"}},
		{"a", '.', []string{"a", ""}},
		{"a.b", '.', []string{"a", "b"}},
		{"a.b.c", '.', []string{"a", "b.c"}},
		{"a/b/c", '.', []string{"a/b/c", ""}},
		{"a/b/c", '/', []string{"a", "b/c"}},
	}

	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			t.Parallel()

			root, child := splitPath(test.path, test.sep)

			require.Equal(t, test.exp[0], root)
			require.Equal(t, test.exp[1], child)
		})
	}
}
