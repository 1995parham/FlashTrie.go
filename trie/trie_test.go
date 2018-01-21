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
	if trie.Root.NextHop == "" || trie.Root.prefix != "" {
		t.Fatal("Invalid Route Insertation: *")
	}

	trie.Add("11*", "B")
	if trie.Root.Right.Right.NextHop == "" || trie.Root.Right.Right.prefix != "11" {
		t.Fatal("Invalid Route Insertation: 11*")
	}

	if trie.Root.Right.prefix != "1" {
		t.Fatal("Invalid Prefix: 1 != %s", trie.Root.Right.prefix)
	}

	if trie.Height != 3 {
		t.Fatalf("Invalid height: 3 != %d", trie.Height)
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

	if trie.Height != 4 {
		t.Fatalf("Invalid height: 4 != %d", trie.Height)
	}

	tries := trie.Divide(3)
	if len(tries) != 2 {
		t.Fatalf("Invalid number of levels in dividation")
	}
	if len(tries[0]) != 1 {
		t.Fatalf("Invalid number of tires in level 0")
	}
	if tries[0][0].Height != 3 {
		t.Fatalf("Invalid height of trie in level 0: 3 != %d", tries[0][0].Height)
	}
	if len(tries[1]) != 1 {
		t.Fatalf("Invalid number of tires in level 1")
	}
}

func TestLookup1(t *testing.T) {
	trie := New()

	trie.Add("*", "A")
	trie.Add("11*", "B")
	trie.Add("10*", "C")

	if trie.Lookup("11101") != "B" {
		t.Fatalf("Invalid lookup for 11101. B != %s", trie.Lookup("11101"))
	}
}

func TestLookup2(t *testing.T) {
	trie := New()

	trie.Add("*", "A")
	trie.Add("1*", "B")
	trie.Add("00*", "C")
	trie.Add("11*", "D")
	trie.Add("100*", "E")

	if trie.Lookup("100") != "E" {
		t.Fatalf("Invalid lookup for 100. E != %s", trie.Lookup("100"))
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
