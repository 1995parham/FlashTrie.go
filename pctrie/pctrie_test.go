package pctrie_test

import (
	"testing"

	"github.com/1995parham/FlashTrie.go/pctrie"
	"github.com/1995parham/FlashTrie.go/trie"
)

func TestBasic2(t *testing.T) {
	t.Parallel()

	trie := trie.New()

	trie.Add("*", "A")
	trie.Add("1*", "B")
	trie.Add("00*", "C")
	trie.Add("11*", "D")
	trie.Add("100*", "E")

	pctrie := pctrie.New(trie, 4)
	if string(pctrie.Bitmap) != "101" {
		t.Fatalf("Invalid bitmap. 101 != %s", string(pctrie.Bitmap))
	}

	for i, b := range pctrie.Bitmap {
		if b == '1' {
			if len(pctrie.NextHops[i]) == 0 {
				t.Fatalf("Invalid NextHops at %d\n", i)
			}
		} else {
			if len(pctrie.NextHops[i]) != 0 {
				t.Fatalf("Invalid NextHops at %d\n", i)
			}
		}
	}
}
