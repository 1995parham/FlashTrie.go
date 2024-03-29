package pctrie

import (
	"testing"

	"github.com/1995parham/FlashTrie.go/trie"
)

func TestBasic1(t *testing.T) {
	t.Parallel()

	trie := trie.New()

	trie.Add("*", "A")
	trie.Add("1*", "B")
	trie.Add("00*", "C")
	trie.Add("11*", "D")
	trie.Add("100*", "E")

	pctrie := New(trie, 2)
	if pctrie.Size != 7 {
		t.Fatalf("Invalid bitmap size. 7 != %d", pctrie.Size)
	}

	if pctrie.compBits != 1 {
		t.Fatalf("Invalid number of bits are used to identify NHI. 1 != %d", pctrie.compBits)
	}

	if string(pctrie.Bitmap) != "0110010" {
		t.Fatalf("Invalid bitmap. 0110010 != %s", string(pctrie.Bitmap))
	}

	if pctrie.Lookup("100") != trie.Lookup("100") {
		t.Fatalf("Invalid lookup for 100. %s != %s", trie.Lookup("100"), pctrie.Lookup("100"))
	}

	if pctrie.Lookup("001") != trie.Lookup("001") {
		t.Fatalf("Invalid lookup for 100. %s != %s", trie.Lookup("001"), pctrie.Lookup("001"))
	}
}
