package geohash

import (
	"errors"
	"fmt"

	"github.com/1995parham/FlashTrie.go/fltrie"
)

const (
	base32Alphabet = "0123456789bcdefghjkmnpqrstuvwxyz"

	bitsPerChar      = 5
	defaultPrecision = 12
	defaultStride    = 5
	compSize         = 2
	trieDepth        = 2

	maxLat = 90.0
	maxLng = 180.0
)

var (
	ErrOutOfRange       = errors.New("coordinate out of range")
	ErrInvalidCharacter = errors.New("invalid geohash character")
	ErrEmptyHash        = errors.New("empty geohash string")
	ErrZeroPrecision    = errors.New("precision must be > 0")
)

// charToIndex maps each base-32 character to its 5-bit index.
// Initialized once at package load; effectively immutable after that.
//
//nolint:gochecknoglobals
var charToIndex = buildCharIndex()

func buildCharIndex() [256]int {
	var table [256]int

	for i := range table {
		table[i] = -1
	}

	for i, c := range base32Alphabet {
		table[c] = i
	}

	return table
}

// Coord represents a geographic coordinate.
type Coord struct {
	Lat float64 // [-90, 90]
	Lng float64 // [-180, 180]
}

// Adapter implements fltrie.Adapter[Coord] with configurable precision.
type Adapter struct {
	precision uint // number of geohash characters (each = 5 bits)
}

// NewAdapter creates a new geohash adapter with the given precision (1..12).
func NewAdapter(precision uint) Adapter {
	return Adapter{precision: precision}
}

// encodeBit performs one step of the geohash binary subdivision.
// It returns '1' and narrows to the upper half when val >= mid,
// or '0' and narrows to the lower half otherwise.
func encodeBit(val, lo, hi float64) (byte, float64, float64) {
	mid := (lo + hi) / 2 //nolint:mnd

	if val >= mid {
		return '1', mid, hi
	}

	return '0', lo, mid
}

// Encode converts a coordinate to a binary string of length precision*5.
func (a Adapter) Encode(c Coord) (string, error) {
	if a.precision == 0 {
		return "", ErrZeroPrecision
	}

	if c.Lat < -maxLat || c.Lat > maxLat || c.Lng < -maxLng || c.Lng > maxLng {
		return "", fmt.Errorf("%w: lat=%f, lng=%f", ErrOutOfRange, c.Lat, c.Lng)
	}

	totalBits := a.precision * bitsPerChar
	bits := make([]byte, totalBits)

	lngMin, lngMax := -maxLng, maxLng
	latMin, latMax := -maxLat, maxLat

	for i := range totalBits {
		if i%2 == 0 {
			bits[i], lngMin, lngMax = encodeBit(c.Lng, lngMin, lngMax)
		} else {
			bits[i], latMin, latMax = encodeBit(c.Lat, latMin, latMax)
		}
	}

	return string(bits), nil
}

// KeyBits returns the total number of bits for this adapter.
func (a Adapter) KeyBits() uint {
	return a.precision * bitsPerChar
}

// ParseGeohash converts a geohash string to a binary prefix string for fltrie.Add.
func ParseGeohash(hash string) (string, error) {
	if len(hash) == 0 {
		return "", ErrEmptyHash
	}

	bits := make([]byte, 0, len(hash)*bitsPerChar)

	for i := range len(hash) {
		idx := charToIndex[hash[i]]
		if idx < 0 {
			return "", fmt.Errorf("%w: %q", ErrInvalidCharacter, hash[i])
		}

		bits = append(bits,
			'0'+byte((idx>>4)&1), //nolint:mnd
			'0'+byte((idx>>3)&1), //nolint:mnd
			'0'+byte((idx>>2)&1), //nolint:mnd
			'0'+byte((idx>>1)&1),
			'0'+byte(idx&1),
		)
	}

	return string(bits) + "*", nil
}

// DefaultConfig returns the standard FlashTrie config for geohash lookups.
// KeyBits=60 (12 chars x 5 bits), Stride=5 (one geohash char).
func DefaultConfig() fltrie.Config {
	return fltrie.Config{
		KeyBits:   defaultPrecision * bitsPerChar,
		Stride:    defaultStride,
		CompSize:  compSize,
		TrieDepth: trieDepth,
	}
}
