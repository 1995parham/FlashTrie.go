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
	"fmt"

	"github.com/1995parham/FlashTrie.go/trie"
)

// PCTrie represents prefix-compresed trie data structure
type PCTrie struct {
	Bitmap   int
	NextHops [][]string
	Size     int
}

// New creates new prefix-compresed trie
func New(tr *trie.Trie, compSize int) *PCTrie {
	nodes := tr.ToArray()
	nhs := make([][]string, 0)
	bitmap := 0
	size := 0

	for s := compSize; s < len(nodes)-1; s += compSize {
		empty := true
		nh := make([]string, compSize)

		for t := s; t < s+compSize; t++ {
			nh[t-s] = tr.Lookup(fmt.Sprintf("%b", t)[1:] + "*")
			if nodes[t].NextHop != "" {
				empty = false
			}
		}
		if !empty {
			bitmap |= 1 << uint(size)
			nhs = append(nhs, nh)
		} else {
			nhs = append(nhs, make([]string, 0))
		}
		size++
	}

	// Eliminate Redundancy
	for i := 0; i <= size/2; i++ {
		if bitmap&(1<<uint(i)) != 0 {
			if bitmap&(1<<uint(2*i+1)) != 0 && bitmap&(1<<uint(2*i+2)) != 0 {
				bitmap ^= 1 << uint(i)
				nhs[i] = make([]string, 0)
			}
		}
	}

	return &PCTrie{
		Bitmap:   bitmap,
		Size:     size,
		NextHops: nhs,
	}
}
