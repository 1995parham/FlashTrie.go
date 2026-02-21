package trie_test

import (
	"testing"

	"github.com/1995parham/FlashTrie.go/trie"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDivide1(t *testing.T) {
	t.Parallel()
	//          A
	//       /     \
	//      -       B
	//     / \    /  \
	//    C  -   -    D
	//   / \/ \ / \  / \
	//  -  -  - E - -  -

	tr := trie.New[string]()

	tr.Add("*", "A")
	tr.Add("1*", "B")
	tr.Add("00*", "C")
	tr.Add("11*", "D")
	tr.Add("100*", "E")

	require.Equal(t, uint(4), tr.Height)

	tries := tr.Divide(3)
	require.Len(t, tries, 2, "Invalid number of levels in division")
	require.Len(t, tries[0], 1, "Invalid number of tries in level 0")
	assert.Equal(t, uint(3), tries[0][0].Height)
	require.Len(t, tries[1], 1, "Invalid number of tries in level 1")
}

func TestDivide2(t *testing.T) {
	t.Parallel()
	//          A
	//       /     \
	//      -       B
	//     / \    /  \
	//    C  -   -    D
	//   / \/ \ / \  / \
	//  -  -  - E - -  -

	tr := trie.New[string]()

	tr.Add("*", "A")
	tr.Add("1*", "B")
	tr.Add("00*", "C")
	tr.Add("11*", "D")
	tr.Add("100*", "E")

	require.Equal(t, uint(4), tr.Height)

	tries := tr.Divide(2)
	require.Len(t, tries, 2, "Invalid number of levels in division")
	require.Len(t, tries[0], 1, "Invalid number of tries in level 0")
	assert.Equal(t, uint(2), tries[0][0].Height)
	require.Len(t, tries[1], 3, "Invalid number of tries in level 1")

	for _, tr := range tries[1] {
		t.Log(tr.Prefix)
		require.NotNil(t, tr.Root.Value, "Subtries must be independent")
	}
}

func TestLookup1(t *testing.T) {
	t.Parallel()

	tr := trie.New[string]()

	tr.Add("*", "A")
	tr.Add("11*", "B")
	tr.Add("10*", "C")

	result, found := tr.Lookup("11101")
	require.True(t, found)
	assert.Equal(t, "B", result)
}

func TestLookup2(t *testing.T) {
	t.Parallel()

	tr := trie.New[string]()

	tr.Add("*", "A")
	tr.Add("1*", "B")
	tr.Add("00*", "C")
	tr.Add("11*", "D")
	tr.Add("100*", "E")

	result, found := tr.Lookup("100")
	require.True(t, found)
	assert.Equal(t, "E", result)
}

func TestArray1(t *testing.T) {
	t.Parallel()

	tr := trie.New[string]()

	tr.Add("*", "A")
	tr.Add("1*", "B")
	tr.Add("0*", "C")

	nodes := tr.ToArray()

	require.NotNil(t, nodes[1].Value)
	assert.Equal(t, "A", *nodes[1].Value)
	require.NotNil(t, nodes[3].Value)
	assert.Equal(t, "B", *nodes[3].Value)
	require.NotNil(t, nodes[2].Value)
	assert.Equal(t, "C", *nodes[2].Value)
}

func TestAll(t *testing.T) {
	t.Parallel()

	tr := trie.New[string]()

	tr.Add("*", "A")
	tr.Add("1*", "B")
	tr.Add("00*", "C")
	tr.Add("11*", "D")
	tr.Add("100*", "E")

	got := make(map[string]string)
	for prefix, value := range tr.All() {
		got[prefix] = value
	}

	assert.Len(t, got, 5)
	assert.Equal(t, "A", got[""])
	assert.Equal(t, "B", got["1"])
	assert.Equal(t, "C", got["00"])
	assert.Equal(t, "D", got["11"])
	assert.Equal(t, "E", got["100"])
}

func TestAllEarlyBreak(t *testing.T) {
	t.Parallel()

	tr := trie.New[string]()
	tr.Add("*", "A")
	tr.Add("1*", "B")
	tr.Add("0*", "C")

	count := 0
	for range tr.All() {
		count++
		if count == 2 {
			break
		}
	}

	assert.Equal(t, 2, count)
}

func TestMatches(t *testing.T) {
	t.Parallel()

	tr := trie.New[string]()
	tr.Add("*", "A")
	tr.Add("1*", "B")
	tr.Add("10*", "C")
	tr.Add("100*", "E")

	var prefixes []string
	var values []string

	for prefix, value := range tr.Matches("1001") {
		prefixes = append(prefixes, prefix)
		values = append(values, value)
	}

	require.Len(t, prefixes, 4)
	assert.Equal(t, []string{"", "1", "10", "100"}, prefixes)
	assert.Equal(t, []string{"A", "B", "C", "E"}, values)
}

func TestMatchesPartial(t *testing.T) {
	t.Parallel()

	tr := trie.New[string]()
	tr.Add("*", "A")
	tr.Add("11*", "B")

	var values []string

	for _, value := range tr.Matches("10") {
		values = append(values, value)
	}

	// Only root matches; "11*" does not match route "10"
	require.Len(t, values, 1)
	assert.Equal(t, "A", values[0])
}

func BenchmarkAdd(b *testing.B) {
	for b.Loop() {
		tr := trie.New[string]()
		tr.Add("*", "A")
		tr.Add("1*", "B")
		tr.Add("00*", "C")
		tr.Add("11*", "D")
		tr.Add("100*", "E")
	}
}

func BenchmarkLookup(b *testing.B) {
	tr := trie.New[string]()
	tr.Add("*", "A")
	tr.Add("1*", "B")
	tr.Add("00*", "C")
	tr.Add("11*", "D")
	tr.Add("100*", "E")

	b.ResetTimer()

	for b.Loop() {
		tr.Lookup("100")
	}
}
