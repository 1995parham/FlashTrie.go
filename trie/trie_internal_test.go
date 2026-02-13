package trie

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdd1(t *testing.T) {
	t.Parallel()

	trie := New[string]()

	trie.Add("*", "A")

	require.NotNil(t, trie.Root.Value, "Invalid Route Insertion: *")
	assert.Empty(t, trie.Root.prefix)

	trie.Add("11*", "B")

	require.NotNil(t, trie.Root.Right.Right.Value, "Invalid Route Insertion: 11*")
	assert.Equal(t, "11", trie.Root.Right.Right.prefix)
	assert.Equal(t, "1", trie.Root.Right.prefix)
	assert.Equal(t, uint(3), trie.Height)
}

func TestEmptyTrieLookup(t *testing.T) {
	t.Parallel()

	trie := New[string]()
	_, found := trie.Lookup("100")
	assert.False(t, found)
}

func TestDuplicateRoutes(t *testing.T) {
	t.Parallel()

	trie := New[string]()
	trie.Add("*", "A")
	trie.Add("*", "B")
	result, found := trie.Lookup("100")
	require.True(t, found)
	assert.Equal(t, "B", result)
}

func TestToArrayPanicsOnLargeHeight(t *testing.T) {
	t.Parallel()

	trie := New[string]()
	trie.Height = 25

	assert.Panics(t, func() {
		trie.ToArray()
	})
}
