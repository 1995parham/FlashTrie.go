package fltrie

import (
	"errors"
	"fmt"

	"github.com/1995parham/FlashTrie.go/pctrie"
	"github.com/1995parham/FlashTrie.go/trie"
)

var (
	ErrAlreadyBuilt  = errors.New("flash trie is already built")
	ErrInvalidHeight = errors.New("flash trie height does not match config KeyBits")
	ErrInvalidRoute  = errors.New("encoded route length does not match config KeyBits")
)

// Adapter encodes a query key into a full-length binary string for lookup.
type Adapter[K any] interface {
	Encode(key K) (string, error)
	KeyBits() uint
}

// Config holds the parameters for a FlashTrie.
type Config struct {
	KeyBits   uint // total key length in bits (32 for IPv4, 128 for IPv6, etc.)
	Stride    uint // bits per subdivision level (default: 8)
	CompSize  int  // pctrie compression factor (default: 2)
	TrieDepth uint // number of stride-levels covered by direct trie (default: 2)
}

type hashElement[V any] struct {
	pctrie *pctrie.PCTrie[V]
	value  *V // root value of the subtrie
}

// FLTrie represents flash trie structure.
type FLTrie[K, V any] struct {
	adapter Adapter[K]
	config  Config
	trie    *trie.Trie[V]
	pctries []map[string]*hashElement[V]
	built   bool
}

// New creates empty and unbuilt flash trie.
func New[K, V any](adapter Adapter[K], cfg Config) *FLTrie[K, V] {
	// Calculate number of pctrie levels
	numDivideLevels := cfg.KeyBits / cfg.Stride
	if cfg.KeyBits%cfg.Stride != 0 {
		numDivideLevels++
	}

	numPCTrieLevels := numDivideLevels - cfg.TrieDepth

	return &FLTrie[K, V]{
		adapter: adapter,
		config:  cfg,
		built:   false,
		trie:    trie.New[V](),
		pctries: make([]map[string]*hashElement[V], numPCTrieLevels),
	}
}

// Add adds new route into unbuilt flash trie.
// Given route must be in binary regex format e.g. *, 11*.
func (fl *FLTrie[K, V]) Add(route string, value V) error {
	if fl.built {
		return ErrAlreadyBuilt
	}

	fl.trie.Add(route, value)

	return nil
}

// Build builds flash trie multi-level hierarchy.
func (fl *FLTrie[K, V]) Build() error {
	if fl.trie.Height != fl.config.KeyBits {
		return fmt.Errorf("%w: got %d, want %d", ErrInvalidHeight, fl.trie.Height, fl.config.KeyBits)
	}

	if fl.built {
		return ErrAlreadyBuilt
	}

	fl.built = true

	tries := fl.trie.Divide(fl.config.Stride)

	numLevels := uint(len(tries))
	for lvl := fl.config.TrieDepth; lvl < numLevels; lvl++ {
		idx := lvl - fl.config.TrieDepth

		fl.pctries[idx] = make(map[string]*hashElement[V])
		for _, sub := range tries[lvl] {
			fl.pctries[idx][sub.Prefix] = &hashElement[V]{
				pctrie: pctrie.New[V](sub, fl.config.CompSize),
				value:  sub.Root.Value,
			}
		}
	}

	return nil
}

// Lookup looks up given key and returns found value.
func (fl *FLTrie[K, V]) Lookup(key K) (V, bool, error) {
	var zero V

	if !fl.built {
		return zero, false, nil
	}

	route, err := fl.adapter.Encode(key)
	if err != nil {
		return zero, false, fmt.Errorf("adapter encode: %w", err)
	}

	if uint(len(route)) != fl.config.KeyBits {
		return zero, false, fmt.Errorf("%w: got %d, want %d", ErrInvalidRoute, len(route), fl.config.KeyBits)
	}

	firstLevelBits := fl.config.TrieDepth * fl.config.Stride

	var result V

	hasResult := false

	// Level 1 (direct trie)
	if val, found := fl.trie.Lookup(route[:firstLevelBits]); found {
		result = val
		hasResult = true
	}

	// Dynamic pctrie levels
	prefixEnd := firstLevelBits

	for i := range fl.pctries {
		if he, ok := fl.pctries[i][route[:prefixEnd]]; ok {
			nhi, found, lookupErr := he.pctrie.Lookup(route[prefixEnd : prefixEnd+fl.config.Stride])
			if lookupErr != nil {
				return zero, false, fmt.Errorf("level %d lookup: %w", i+int(fl.config.TrieDepth), lookupErr)
			}

			if found {
				result, hasResult = nhi, true
			} else if he.value != nil {
				result, hasResult = *he.value, true
			}
		}

		prefixEnd += fl.config.Stride
	}

	return result, hasResult, nil
}
