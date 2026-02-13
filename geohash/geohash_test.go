package geohash_test

import (
	"strings"
	"testing"

	"github.com/1995parham/FlashTrie.go/fltrie"
	"github.com/1995parham/FlashTrie.go/geohash"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseGeohash(t *testing.T) {
	t.Parallel()

	result, err := geohash.ParseGeohash("u4pru")
	require.NoError(t, err)
	// u=11010, 4=00100, p=10101, r=10111, u=11010
	assert.Equal(t, "1101000100101011011111010*", result)
	assert.Len(t, result, 26) // 25 bits + "*"
}

func TestParseGeohashSingleChar(t *testing.T) {
	t.Parallel()

	result, err := geohash.ParseGeohash("u")
	require.NoError(t, err)
	assert.Equal(t, "11010*", result)
	assert.Len(t, result, 6) // 5 bits + "*"
}

func TestParseGeohashInvalidCharA(t *testing.T) {
	t.Parallel()

	_, err := geohash.ParseGeohash("a")
	require.Error(t, err)
}

func TestParseGeohashInvalidCharI(t *testing.T) {
	t.Parallel()

	_, err := geohash.ParseGeohash("i")
	require.Error(t, err)
}

func TestParseGeohashInvalidCharL(t *testing.T) {
	t.Parallel()

	_, err := geohash.ParseGeohash("l")
	require.Error(t, err)
}

func TestParseGeohashInvalidCharO(t *testing.T) {
	t.Parallel()

	_, err := geohash.ParseGeohash("o")
	require.Error(t, err)
}

func TestParseGeohashEmpty(t *testing.T) {
	t.Parallel()

	_, err := geohash.ParseGeohash("")
	require.Error(t, err)
}

func TestAdapterEncode(t *testing.T) {
	t.Parallel()

	adapter := geohash.NewAdapter(12)
	result, err := adapter.Encode(geohash.Coord{Lat: 57.64911, Lng: 10.40744})
	require.NoError(t, err)
	assert.Len(t, result, 60)
}

func TestAdapterEncodeKnownCoord(t *testing.T) {
	t.Parallel()

	// lat=57.64911, lng=10.40744 should produce geohash starting with "u4pru"
	// The first 25 bits should match ParseGeohash("u4pru") without the "*"
	adapter := geohash.NewAdapter(5)
	result, err := adapter.Encode(geohash.Coord{Lat: 57.64911, Lng: 10.40744})
	require.NoError(t, err)
	assert.Equal(t, "1101000100101011011111010", result)
}

func TestAdapterEncodeOutOfRange(t *testing.T) {
	t.Parallel()

	adapter := geohash.NewAdapter(5)

	_, err := adapter.Encode(geohash.Coord{Lat: 91, Lng: 0})
	require.Error(t, err)

	_, err = adapter.Encode(geohash.Coord{Lat: 0, Lng: 181})
	require.Error(t, err)

	_, err = adapter.Encode(geohash.Coord{Lat: -91, Lng: 0})
	require.Error(t, err)

	_, err = adapter.Encode(geohash.Coord{Lat: 0, Lng: -181})
	require.Error(t, err)
}

func TestAdapterKeyBits(t *testing.T) {
	t.Parallel()

	adapter := geohash.NewAdapter(12)
	assert.Equal(t, uint(60), adapter.KeyBits())

	adapter5 := geohash.NewAdapter(5)
	assert.Equal(t, uint(25), adapter5.KeyBits())
}

func TestRoundTrip(t *testing.T) {
	t.Parallel()

	// Use precision=6 (30-bit keys) so we can add shorter geohash prefixes
	// and still have room for the trie depth.
	adapter := geohash.NewAdapter(6)
	cfg := fltrie.Config{
		KeyBits:   30,
		Stride:    5,
		CompSize:  2,
		TrieDepth: 2,
	}

	fl := fltrie.New[geohash.Coord, string](adapter, cfg)

	// Add a 4-char geohash "u4pr" as a route (20 bits).
	prefix, err := geohash.ParseGeohash("u4pr")
	require.NoError(t, err)
	require.NoError(t, fl.Add(prefix, "Aalborg"))

	// The trie needs Height == KeyBits (30). The longest prefix must be 29 bits.
	// Add a dummy 29-bit route in a non-overlapping region to satisfy this.
	require.NoError(t, fl.Add(strings.Repeat("0", 29)+"*", "dummy"))

	require.NoError(t, fl.Build())

	val, found, err := fl.Lookup(geohash.Coord{Lat: 57.64911, Lng: 10.40744})
	require.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, "Aalborg", val)
}

func TestEncodeParseConsistency(t *testing.T) {
	t.Parallel()

	// Verify that Encode and ParseGeohash produce compatible binary strings.
	adapter := geohash.NewAdapter(5)
	encoded, err := adapter.Encode(geohash.Coord{Lat: 57.64911, Lng: 10.40744})
	require.NoError(t, err)

	prefix, err := geohash.ParseGeohash("u4pru")
	require.NoError(t, err)

	assert.Equal(t, encoded, strings.TrimSuffix(prefix, "*"))
}

func BenchmarkParseGeohash(b *testing.B) {
	for b.Loop() {
		_, _ = geohash.ParseGeohash("u4pruydqqvj8")
	}
}

func BenchmarkAdapterEncode(b *testing.B) {
	adapter := geohash.NewAdapter(12)
	coord := geohash.Coord{Lat: 57.64911, Lng: 10.40744}

	b.ResetTimer()

	for b.Loop() {
		_, _ = adapter.Encode(coord)
	}
}
