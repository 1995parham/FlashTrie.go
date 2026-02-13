package pctrie_test

import (
	"testing"

	"github.com/1995parham/FlashTrie.go/pctrie"
	"github.com/1995parham/FlashTrie.go/trie"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBasic2(t *testing.T) {
	t.Parallel()

	tr := trie.New[string]()

	tr.Add("*", "A")
	tr.Add("1*", "B")
	tr.Add("00*", "C")
	tr.Add("11*", "D")
	tr.Add("100*", "E")

	pc := pctrie.New[string](tr, 4)
	assert.Equal(t, "101", pc.Bitmap.String())

	for i := range pc.Bitmap.Len() {
		if pc.Bitmap.Get(i) {
			require.NotEmpty(t, pc.NextHops[i], "Invalid NextHops at %d", i)
		} else {
			require.Empty(t, pc.NextHops[i], "Invalid NextHops at %d", i)
		}
	}
}
