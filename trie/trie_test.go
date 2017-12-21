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

import (
	"testing"
)

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

	if trie.Stride != 3 {
		t.Fatalf("Invalid stride: 3 != %d", trie.Stride)
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

func TestArray1(t *testing.T) {
	trie := New()

	trie.Add("*", "A")
	trie.Add("1*", "B")
	trie.Add("0*", "C")

	nodes := trie.ToArray()

	if nodes[1].NextHop != "A" {
		t.Fatal("Invalid Node Array")
	}
	if nodes[2].NextHop != "B" {
		t.Fatal("Invalid Node Array")
	}
	if nodes[3].NextHop != "C" {
		t.Fatal("Invalid Node Array")
	}
}
