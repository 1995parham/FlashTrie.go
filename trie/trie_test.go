package trie_test

import (
	"testing"

	"github.com/1995parham/FlashTrie.go/trie"
)

func TestDivide1(t *testing.T) {
	t.Parallel()
	//          A
	//       /     \
	//      -       B
	//     / \    /  \
	//    C  -   -    D
	//   / \/ \ / \  / \
	//  -  -  - E - -  -

	trie := trie.New()

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
		t.Fatalf("Invalid number of tires in level 0: 1 != %d", len(tries[0]))
	}

	if tries[0][0].Height != 3 {
		t.Fatalf("Invalid height of trie in level 0: 3 != %d", tries[0][0].Height)
	}

	if len(tries[1]) != 1 {
		t.Fatalf("Invalid number of tires in level 1: 1 != %d", len(tries[1]))
	}
}

func TestDivide2(t *testing.T) {
	t.Parallel()
	//          A
	//       /     \
	//      -       B
	//     / \    /  \
	//    C  -   -    D
	//   / \/ \ / \  / \
	//  -  -  - E - -  -

	trie := trie.New()

	trie.Add("*", "A")
	trie.Add("1*", "B")
	trie.Add("00*", "C")
	trie.Add("11*", "D")
	trie.Add("100*", "E")

	if trie.Height != 4 {
		t.Fatalf("Invalid height: 4 != %d", trie.Height)
	}

	tries := trie.Divide(2)
	if len(tries) != 2 {
		t.Fatalf("Invalid number of levels in dividation")
	}

	if len(tries[0]) != 1 {
		t.Fatalf("Invalid number of tires in level 0: 1 != %d", len(tries[0]))
	}

	if tries[0][0].Height != 2 {
		t.Fatalf("Invalid height of trie in level 0: 2 != %d", tries[0][0].Height)
	}

	if len(tries[1]) != 3 {
		t.Fatalf("Invalid number of tires in level 1: 3 != %d", len(tries[1]))
	}

	for _, trie := range tries[1] {
		t.Log(trie.Prefix)

		if trie.Root.NextHop == "" {
			t.Fatalf("Subtries must be independent")
		}
	}
}

func TestLookup1(t *testing.T) {
	t.Parallel()

	trie := trie.New()

	trie.Add("*", "A")
	trie.Add("11*", "B")
	trie.Add("10*", "C")

	if trie.Lookup("11101") != "B" {
		t.Fatalf("Invalid lookup for 11101. B != %s", trie.Lookup("11101"))
	}
}

func TestLookup2(t *testing.T) {
	t.Parallel()

	trie := trie.New()

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
	t.Parallel()

	trie := trie.New()

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
