package filter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewStringMatcherNil(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	sm := NewStringMatcher(nil, '/')
	require.Nil(sm.node.children)
	require.Equal(noval{}, sm.value)
	require.False(sm.set)
}

func TestNewStringMatcherEmpty(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	keys := []string{}

	sm := NewStringMatcher(keys, '/')
	require.Nil(sm.node.children)
	require.Equal(noval{}, sm.value)
	require.False(sm.set)
}

func TestNewStringMatcher(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	keys := []string{"a"}

	sm := NewStringMatcher(keys, '/')
	require.NotNil(sm.node.children)
	require.Equal(noval{}, sm.value)
	require.False(sm.set)
}

func TestStringMatcherMatch(t *testing.T) {
	t.Parallel()

	sm := NewStringMatcher(nil, '/')

	ok := sm.Match("")
	require.False(t, ok)
}

func TestStringMatcher(t *testing.T) {
	t.Parallel()

	sm := NewStringMatcher([]string{"a", "x/y", "/"}, '/')

	tests := []struct {
		path string
		exp  bool
	}{
		{"a", true},
		{"a/b", true},
		{"x/y", true},
		{"x/y/z", true},
		{"x/v", false},
		{"/", true},
		{"/a", false},
	}

	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			t.Parallel()

			val := sm.Match(test.path)
			require.Equal(t, test.exp, val)
		})
	}
}
