package pctrie

import (
	"fmt"
	"strconv"

	"github.com/1995parham/FlashTrie.go/trie"
)

// PCTrie represents prefix-compressed trie data structure.
type PCTrie[V any] struct {
	Bitmap   *Bitmap
	NextHops [][]*V
	Size     int // Bitmap size

	height   uint // number of bits for routing
	compBits uint // number of bits are used to identify the corresponding NHI
}

// New creates new prefix-compressed trie.
// nolint: cyclop, funlen
func New[V any](tr *trie.Trie[V], compSize int) *PCTrie[V] {
	nodes := tr.ToArray()
	nhs := make([][]*V, 0)

	// Count entries first
	size := 0
	for s := compSize; s < len(nodes)-1; s += compSize {
		size++
	}

	bm := newBitmap(size)
	idx := 0

	// All PC-Trie node
	for s := compSize; s < len(nodes)-1; s += compSize {
		empty := true
		nh := make([]*V, compSize)

		for t := s; t < s+compSize; t++ {
			val, found := tr.Lookup(fmt.Sprintf("%b", t)[1:])
			if found {
				nh[t-s] = &val
			}

			if nodes[t].Value != nil {
				empty = false
			}
		}

		if !empty {
			bm.Set(idx)

			nhs = append(nhs, nh)
		} else {
			nhs = append(nhs, make([]*V, 0))
		}

		idx++
	}

	// Eliminate Redundancy
	for i := range size / 2 {
		if bm.Get(i) {
			// nolint: mnd
			if bm.Get(2*i+1) && bm.Get(2*i+2) {
				bm.Clear(i)

				nhs[i] = make([]*V, 0)
			}
		}
	}

	// Calculate number of bits are needed to identify NHI
	var compBits uint
	for compSize>>compBits != 1 {
		compBits++
	}

	return &PCTrie[V]{
		Bitmap:   bm,
		Size:     size,
		NextHops: nhs,

		height:   tr.Height - 1,
		compBits: compBits,
	}
}

// Lookup looks up given route in pc-trie and returns found value.
// Given route must be in binary representation e.g. 111111..
// Note that this function assumes that given route length is greater than
// original trie height.
func (pc *PCTrie[V]) Lookup(route string) (V, bool, error) {
	var zero V

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
		i, err := strconv.ParseInt(b[bits-pc.compBits:], 2, 0)
		if err != nil {
			return zero, false, fmt.Errorf("invalid route binary at NHI: %w", err)
		}

		nhi = int(i)

		// bitmap access
		i, err = strconv.ParseInt(b[:bits-pc.compBits], 2, 0)
		if err != nil {
			return zero, false, fmt.Errorf("invalid route binary at offset: %w", err)
		}

		// nolint: unconvert
		offset = int(i) + (1 << uint(bits-pc.compBits)) - 1

		if pc.Bitmap.Get(offset) {
			indicator = true
		}

		bits--
	}

	if !indicator {
		return zero, false, nil
	}

	if nhi < len(pc.NextHops[offset]) && pc.NextHops[offset][nhi] != nil {
		return *pc.NextHops[offset][nhi], true, nil
	}

	return zero, false, nil
}
