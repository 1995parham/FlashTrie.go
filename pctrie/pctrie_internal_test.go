package pctrie

import (
	"slices"
	"testing"

	"github.com/1995parham/FlashTrie.go/trie"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBasic1(t *testing.T) {
	t.Parallel()

	tr := trie.New[string]()

	tr.Add("*", "A")
	tr.Add("1*", "B")
	tr.Add("00*", "C")
	tr.Add("11*", "D")
	tr.Add("100*", "E")

	pc := New[string](tr, 2)
	assert.Equal(t, 7, pc.Size)
	assert.Equal(t, uint(1), pc.compBits)
	assert.Equal(t, "0110010", pc.Bitmap.String())

	result, found, err := pc.Lookup("100")
	require.NoError(t, err)
	require.True(t, found)

	trResult, trFound := tr.Lookup("100")
	require.True(t, trFound)
	assert.Equal(t, trResult, result)

	result, found, err = pc.Lookup("001")
	require.NoError(t, err)
	require.True(t, found)

	trResult, trFound = tr.Lookup("001")
	require.True(t, trFound)
	assert.Equal(t, trResult, result)
}

func TestBitmapOnes(t *testing.T) {
	t.Parallel()

	bm := newBitmap(10)
	bm.Set(0)
	bm.Set(3)
	bm.Set(7)
	bm.Set(9)

	indices := slices.Collect(bm.Ones())

	assert.Equal(t, []int{0, 3, 7, 9}, indices)
}

func TestBitmapOnesEmpty(t *testing.T) {
	t.Parallel()

	bm := newBitmap(64)

	count := 0
	for range bm.Ones() {
		count++
	}

	assert.Zero(t, count)
}

func TestBitmapOnesEarlyBreak(t *testing.T) {
	t.Parallel()

	bm := newBitmap(10)
	bm.Set(1)
	bm.Set(5)
	bm.Set(8)

	var indices []int
	for idx := range bm.Ones() {
		indices = append(indices, idx)
		if len(indices) == 2 {
			break
		}
	}

	assert.Equal(t, []int{1, 5}, indices)
}

func BenchmarkPCTrieLookup(b *testing.B) {
	tr := trie.New[string]()
	tr.Add("*", "A")
	tr.Add("1*", "B")
	tr.Add("00*", "C")
	tr.Add("11*", "D")
	tr.Add("100*", "E")

	pc := New[string](tr, 2)

	b.ResetTimer()

	for b.Loop() {
		_, _, _ = pc.Lookup("100")
	}
}
