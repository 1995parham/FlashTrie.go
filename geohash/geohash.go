package geohash

import (
	"errors"
	"fmt"

	"github.com/1995parham/FlashTrie.go/fltrie"
)

const base32Alphabet = "0123456789bcdefghjkmnpqrstuvwxyz"

var (
	ErrOutOfRange       = errors.New("coordinate out of range")
	ErrInvalidCharacter = errors.New("invalid geohash character")
	ErrEmptyHash        = errors.New("empty geohash string")
	ErrZeroPrecision    = errors.New("precision must be > 0")
)

// charToIndex maps each base-32 character to its 5-bit index.
var charToIndex [256]int

func init() {
	for i := range charToIndex {
		charToIndex[i] = -1
	}

	for i, c := range base32Alphabet {
		charToIndex[c] = i
	}
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

// Encode converts a coordinate to a binary string of length precision*5.
func (a Adapter) Encode(c Coord) (string, error) {
	if a.precision == 0 {
		return "", ErrZeroPrecision
	}

	if c.Lat < -90 || c.Lat > 90 || c.Lng < -180 || c.Lng > 180 {
		return "", fmt.Errorf("%w: lat=%f, lng=%f", ErrOutOfRange, c.Lat, c.Lng)
	}

	totalBits := a.precision * 5
	bits := make([]byte, totalBits)

	lngMin, lngMax := -180.0, 180.0
	latMin, latMax := -90.0, 90.0

	for i := range totalBits {
		if i%2 == 0 {
			// even bit: longitude
			mid := (lngMin + lngMax) / 2
			if c.Lng >= mid {
				bits[i] = '1'
				lngMin = mid
			} else {
				bits[i] = '0'
				lngMax = mid
			}
		} else {
			// odd bit: latitude
			mid := (latMin + latMax) / 2
			if c.Lat >= mid {
				bits[i] = '1'
				latMin = mid
			} else {
				bits[i] = '0'
				latMax = mid
			}
		}
	}

	return string(bits), nil
}

// KeyBits returns the total number of bits for this adapter.
func (a Adapter) KeyBits() uint {
	return a.precision * 5
}

// ParseGeohash converts a geohash string to a binary prefix string for fltrie.Add.
func ParseGeohash(hash string) (string, error) {
	if len(hash) == 0 {
		return "", ErrEmptyHash
	}

	bits := make([]byte, 0, len(hash)*5)

	for i := range len(hash) {
		idx := charToIndex[hash[i]]
		if idx < 0 {
			return "", fmt.Errorf("%w: %q", ErrInvalidCharacter, hash[i])
		}

		bits = append(bits,
			'0'+byte((idx>>4)&1),
			'0'+byte((idx>>3)&1),
			'0'+byte((idx>>2)&1),
			'0'+byte((idx>>1)&1),
			'0'+byte(idx&1),
		)
	}

	return string(bits) + "*", nil
}

// DefaultConfig returns the standard FlashTrie config for geohash lookups.
// KeyBits=60 (12 chars × 5 bits), Stride=5 (one geohash char).
func DefaultConfig() fltrie.Config {
	return fltrie.Config{
		KeyBits:   60,
		Stride:    5,
		CompSize:  2,
		TrieDepth: 2,
	}
}
