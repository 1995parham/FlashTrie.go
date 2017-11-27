/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 27-11-2017
 * |
 * | File Name:     trie/trie.go
 * +===============================================
 */

package trie

// Trie represents binary trie for route lookup
type Trie struct {
	Root *Node
}

// New creates new trie
func New() Trie {
	return Trie{
		Root: new(Node),
	}
}

// Add adds new route into trie
// given route must be in binary regex format e.g. *, 11*
func (t Trie) Add(route string) {
	it := t.Root
	for _, b := range route {
		if b == '*' {
			it.Prefix = route
		} else {
			if b == '1' {
				if it.Left == nil {
					it.Left = new(Node)
				}
				it = it.Left
			}
			if b == '0' {
				if it.Right == nil {
					it.Right = new(Node)
				}
				it = it.Right
			}
		}
	}
}
