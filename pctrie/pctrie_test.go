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
	"fmt"
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
	bitmap := fmt.Sprintf("%0*b", pctrie.Size, pctrie.Bitmap)
	if bitmap != "0100110" {
		t.Fatal("Invalid bitmap")
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
	bitmap := fmt.Sprintf("%0*b", pctrie.Size, pctrie.Bitmap)
	if bitmap != "101" {
		t.Fatal("Invalid bitmap")
	}
	for i, b := range bitmap {
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
