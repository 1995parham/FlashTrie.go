package trie

import "testing"

func TestAdd1(t *testing.T) {
	t.Parallel()

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
		t.Fatalf("Invalid Prefix: 1 != %s", trie.Root.Right.prefix)
	}

	if trie.Height != 3 {
		t.Fatalf("Invalid height: 3 != %d", trie.Height)
	}
}
