package ipv4

import (
	"errors"
	"fmt"
	"net/netip"

	"github.com/1995parham/FlashTrie.go/fltrie"
)

const (
	keyBits       = 32
	defaultStride = 8
	compSize      = 2
	trieDepth     = 2
)

var (
	ErrNotIPv4   = errors.New("not an IPv4 address")
	ErrInvalidIP = errors.New("invalid IP address")
)

// Adapter implements fltrie.Adapter[netip.Addr] for IPv4 lookups.
type Adapter struct{}

// NewAdapter creates a new IPv4 adapter.
func NewAdapter() Adapter {
	return Adapter{}
}

// Encode converts an IPv4 address to a 32-bit binary string.
func (Adapter) Encode(addr netip.Addr) (string, error) {
	if !addr.Is4() {
		return "", fmt.Errorf("%w: %s", ErrNotIPv4, addr)
	}

	octets := addr.As4()

	var str []byte
	for _, octet := range octets {
		str = fmt.Appendf(str, "%08b", octet)
	}

	return string(str), nil
}

// KeyBits returns 32 for IPv4.
func (Adapter) KeyBits() uint {
	return keyBits
}

// ParseCIDR converts a CIDR string to a binary prefix string for fltrie.Add.
func ParseCIDR(cidr string) (string, error) {
	prefix, err := netip.ParsePrefix(cidr)
	if err != nil {
		return "", fmt.Errorf("cidr is not valid: %w", err)
	}

	addr := prefix.Addr()
	if !addr.Is4() {
		return "", fmt.Errorf("%w: %s", ErrNotIPv4, cidr)
	}

	octets := addr.As4()

	var str []byte
	for _, octet := range octets {
		str = fmt.Appendf(str, "%08b", octet)
	}

	return string(str)[:prefix.Bits()] + "*", nil
}

// DefaultConfig returns the standard FlashTrie config for IPv4.
func DefaultConfig() fltrie.Config {
	return fltrie.Config{
		KeyBits:   keyBits,
		Stride:    defaultStride,
		CompSize:  compSize,
		TrieDepth: trieDepth,
	}
}
