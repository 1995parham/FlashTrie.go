package pctrie

const bitsPerWord = 64

// Bitmap represents a compact bit array using uint64 words.
type Bitmap struct {
	words []uint64
	size  int
}

func newBitmap(size int) *Bitmap {
	nwords := (size + bitsPerWord - 1) / bitsPerWord

	return &Bitmap{
		words: make([]uint64, nwords),
		size:  size,
	}
}

// Set sets the bit at position i.
func (b *Bitmap) Set(i int) {
	b.words[i/bitsPerWord] |= 1 << uint(i%bitsPerWord)
}

// Clear clears the bit at position i.
func (b *Bitmap) Clear(i int) {
	b.words[i/bitsPerWord] &^= 1 << uint(i%bitsPerWord)
}

// Get returns true if the bit at position i is set.
func (b *Bitmap) Get(i int) bool {
	return b.words[i/bitsPerWord]&(1<<uint(i%bitsPerWord)) != 0
}

// Len returns the number of bits in the bitmap.
func (b *Bitmap) Len() int {
	return b.size
}

// String returns a binary string representation of the bitmap.
func (b *Bitmap) String() string {
	s := make([]byte, b.size)

	for i := range b.size {
		if b.Get(i) {
			s[i] = '1'
		} else {
			s[i] = '0'
		}
	}

	return string(s)
}
