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
	"strconv"

	"github.com/AUTProjects/FlashTrie.go/trie"
)

// PCTrie represents prefix-compresed trie data structure
type PCTrie struct {
	Bitmap   int
	NextHops [][]string
	Size     int // Bitmap size

	height   uint // number of bits for routing
	compBits uint // number of bits are used to identify the corresponding NHI

}

// New creates new prefix-compresed trie
func New(tr *trie.Trie, compSize int) *PCTrie {
	nodes := tr.ToArray()
	nhs := make([][]string, 0)
	bitmap := 0
	size := 0

	// All PC-Trie node
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

	// Caculate number of bits are needed to identify NHI
	var compBits uint
	for compSize>>compBits != 1 {
		compBits++
	}

	return &PCTrie{
		Bitmap:   bitmap,
		Size:     size,
		NextHops: nhs,

		height:   tr.Height - 1,
		compBits: compBits,
	}
}

// Lookup lookups given route in pc-tire and returns finded nexhop or -
// given route must be in binary represenation e.g. 111111..
// note that this function assume that given route length is greater than
// orignal trie height
func (pc *PCTrie) Lookup(route string) string {
	// access into NextHops array
	offset := 0
	// NHI indicator
	indicator := false
	// corresponding NHI
	nhi := 0
	// how many bit we need for routing
	bits := pc.height

	for !indicator && bits >= pc.compBits+1 {
		// given bitmap
		b := route[:bits]

		// corresponding NHI
		i, _ := strconv.ParseInt(b[bits-pc.compBits:], 2, 0)
		nhi = int(i)

		// bitmap access
		i, _ = strconv.ParseInt(b[:bits-pc.compBits], 2, 0)
		offset = int(i) + (1 << uint(bits-pc.compBits)) - 1
		if pc.Bitmap&(1<<uint(offset)) != 0 {
			indicator = true
		}

		bits--
	}

	return pc.NextHops[offset][nhi]
}
