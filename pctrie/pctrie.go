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
	Bitmap   int
	NextHops [][]string
	Size     int
}

// New creates new prefix-compresed trie
func New(t *trie.Trie, compSize int) *PCTrie {
	nodes := t.ToArray()
	bitmap := 0
	i := 0

	for s := compSize; s < len(nodes)-1; s += compSize {
		empty := true
		for t := s; t < s+compSize; t++ {
			if nodes[t].NextHop != "" {
				empty = false
			}
		}
		if !empty {
			bitmap |= 1 << uint(i)
		} else {
		}
		i++
	}

	return &PCTrie{
		Bitmap: bitmap,
		Size:   i,
	}
}
