package filter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestForward(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		path []testKey
		head string
		tail []testKey
	}{
		{
			name: "nil",
			path: nil,
			head: "",
			tail: nil,
		},
		{
			name: "empty",
			path: testPath(),
			head: "",
			tail: nil,
		},
		{
			name: "single entry",
			path: testPath("a"),
			head: "a",
			tail: []testKey{},
		},
		{
			name: "two entries",
			path: testPath("a", "b"),
			head: "a",
			tail: testPath("b"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			head, tail := forward(test.path)
			require.Equal(t, test.head, head)
			require.Equal(t, test.tail, tail)
		})
	}
}

func TestBackward(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		path []testKey
		head string
		tail []testKey
	}{
		{
			name: "nil",
			path: nil,
			head: "",
			tail: nil,
		},
		{
			name: "empty",
			path: testPath(),
			head: "",
			tail: nil,
		},
		{
			name: "single entry",
			path: testPath("a"),
			head: "a",
			tail: []testKey{},
		},
		{
			name: "two entries",
			path: testPath("a", "b"),
			head: "b",
			tail: testPath("a"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			head, tail := backward(test.path)
			require.Equal(t, test.head, head)
			require.Equal(t, test.tail, tail)
		})
	}
}

func TestNewBackwardFilterNil(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	bf := NewBackwardFilter[testKey, testVal](nil, '/')
	require.Nil(bf.node.children)
	require.Nil(bf.value)
	require.False(bf.set)
}

func TestNewBackwardFilterEmpty(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	keys := map[string]testKey{}

	bf := NewBackwardFilter[testKey](keys, '/')
	require.Nil(bf.node.children)
	require.Nil(bf.value)
	require.False(bf.set)
}

func TestNewBackwardFilter(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	keys := map[string]testKey{
		"a": nil,
	}

	bf := NewBackwardFilter[testKey](keys, '/')
	require.NotNil(bf.node.children)
	require.Nil(bf.value)
	require.False(bf.set)
}

func TestBackwardFilterGet(t *testing.T) {
	t.Parallel()

	bf := NewBackwardFilter[testKey, testVal](nil, '/')

	val, ok := bf.Get(nil)
	require.False(t, ok)
	require.Nil(t, val)
}
