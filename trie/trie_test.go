/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 27-11-2017
 * |
 * | File Name:     trie/trie_test.go
 * +===============================================
 */

package trie

import "testing"

func TestAdd1(t *testing.T) {
	trie := New()

	trie.Add("*", "A")
	if trie.Root.Prefix != "*" {
		t.Fatal("Invalid Route Insertation: *")
	}

	trie.Add("11*", "B")
	if trie.Root.Left.Left.Prefix != "11*" {
		t.Fatal("Invalid Route Insertation: 11*")
	}
}

func TestLookup1(t *testing.T) {
	trie := New()

	trie.Add("*", "A")
	trie.Add("11*", "B")
	trie.Add("10*", "C")

	if trie.Lookup("11101") != "B" {
		t.Fatal("Invalid Route Lookup: 11101")
	}
}
