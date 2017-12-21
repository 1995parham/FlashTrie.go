/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 16-12-2017
 * |
 * | File Name:     pctrie.go
 * +===============================================
 */

package pctrie

import (
	"github.com/1995parham/FlashTrie.go/trie"
)

// PCTrie represents prefix-compresed trie data structure
type PCTrie struct {
}

// New creates new prefix-compresed trie
func New(t *trie.Trie, compSize int) {
	nodes := t.ToArray()

	for s := compSize - 1; s < len(nodes); s += compSize {
		for t := s; t < s+compSize-1; t++ {
		}
	}
}
