/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 21-12-2017
 * |
 * | File Name:     pctrie_test.go
 * +===============================================
 */

package pctrie

import (
	"testing"

	"github.com/AUTProjects/FlashTrie.go/trie"
)

func TestBasic1(t *testing.T) {
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

func TestBasic2(t *testing.T) {
	trie := trie.New()

	trie.Add("*", "A")
	trie.Add("1*", "B")
	trie.Add("00*", "C")
	trie.Add("11*", "D")
	trie.Add("100*", "E")

	pctrie := New(trie, 4)
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
