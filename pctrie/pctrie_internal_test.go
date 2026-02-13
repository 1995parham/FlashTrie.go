package pctrie

import (
	"testing"

	"github.com/1995parham/FlashTrie.go/trie"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBasic1(t *testing.T) {
	t.Parallel()

	tr := trie.New[string]()

	tr.Add("*", "A")
	tr.Add("1*", "B")
	tr.Add("00*", "C")
	tr.Add("11*", "D")
	tr.Add("100*", "E")

	pc := New[string](tr, 2)
	assert.Equal(t, 7, pc.Size)
	assert.Equal(t, uint(1), pc.compBits)
	assert.Equal(t, "0110010", pc.Bitmap.String())

	result, found, err := pc.Lookup("100")
	require.NoError(t, err)
	require.True(t, found)

	trResult, trFound := tr.Lookup("100")
	require.True(t, trFound)
	assert.Equal(t, trResult, result)

	result, found, err = pc.Lookup("001")
	require.NoError(t, err)
	require.True(t, found)

	trResult, trFound = tr.Lookup("001")
	require.True(t, trFound)
	assert.Equal(t, trResult, result)
}

func BenchmarkPCTrieLookup(b *testing.B) {
	tr := trie.New[string]()
	tr.Add("*", "A")
	tr.Add("1*", "B")
	tr.Add("00*", "C")
	tr.Add("11*", "D")
	tr.Add("100*", "E")

	pc := New[string](tr, 2)

	b.ResetTimer()

	for b.Loop() {
		_, _, _ = pc.Lookup("100")
	}
}
