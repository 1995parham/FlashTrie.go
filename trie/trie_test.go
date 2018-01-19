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
	if trie.Root.Right.Right.Prefix != "11*" {
		t.Fatal("Invalid Route Insertation: 11*")
	}

	if trie.Hight != 3 {
		t.Fatalf("Invalid stride: 3 != %d", trie.Hight)
	}
}

func TestDivide(t *testing.T) {
	//          A
	//       /     \
	//      -       B
	//     / \    /  \
	//    C  -   -    D
	//   / \/ \ / \  / \
	//  -  -  - E - -  -

	trie := New()

	trie.Add("*", "A")
	trie.Add("1*", "B")
	trie.Add("00*", "C")
	trie.Add("11*", "D")
	trie.Add("100*", "E")

	tries := trie.Divide(3)
	if len(tries) != 2 {
		t.Fatalf("Invalid number of levels in dividation")
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
	if nodes[3].NextHop != "B" {
		t.Fatal("Invalid Node Array")
	}
	if nodes[2].NextHop != "C" {
		t.Fatal("Invalid Node Array")
	}
}
