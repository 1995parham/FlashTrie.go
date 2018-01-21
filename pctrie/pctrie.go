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
	Bitmap   []byte
	NextHops [][]string
	Size     int // Bitmap size

	height   uint // number of bits for routing
	compBits uint // number of bits are used to identify the corresponding NHI

}

// New creates new prefix-compresed trie
func New(tr *trie.Trie, compSize int) *PCTrie {
	nodes := tr.ToArray()
	nhs := make([][]string, 0)
	bitmap := make([]byte, 0)
	size := 0

	// All PC-Trie node
	for s := compSize; s < len(nodes)-1; s += compSize {
		empty := true
		nh := make([]string, compSize)

		for t := s; t < s+compSize; t++ {
			nh[t-s] = tr.Lookup(fmt.Sprintf("%b", t)[1:])
			if nodes[t].NextHop != "" {
				empty = false
			}
		}
		if !empty {
			bitmap = append(bitmap, '1')
			nhs = append(nhs, nh)
		} else {
			bitmap = append(bitmap, '0')
			nhs = append(nhs, make([]string, 0))
		}
		size++
	}

	// Eliminate Redundancy
	for i := 0; i < size/2; i++ {
		if bitmap[i] == '1' {
			if bitmap[2*i+1] == '1' && bitmap[2*i+2] == '1' {
				bitmap[i] = '0'
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

	for !indicator && bits > pc.compBits {
		// given bitmap
		b := route[:bits]

		// corresponding NHI
		i, _ := strconv.ParseInt(b[bits-pc.compBits:], 2, 0)
		nhi = int(i)

		// bitmap access
		i, _ = strconv.ParseInt(b[:bits-pc.compBits], 2, 0)
		offset = int(i) + (1 << uint(bits-pc.compBits)) - 1
		if pc.Bitmap[offset] == '1' {
			indicator = true
		}

		bits--
	}

	if !indicator {
		return "-"
	}

	return pc.NextHops[offset][nhi]
}
