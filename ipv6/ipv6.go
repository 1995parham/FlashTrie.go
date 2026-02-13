package ipv6

import (
	"errors"
	"fmt"
	"net/netip"

	"github.com/1995parham/FlashTrie.go/fltrie"
)

var (
	ErrNotIPv6   = errors.New("not an IPv6 address")
	ErrInvalidIP = errors.New("invalid IP address")
)

// Adapter implements fltrie.Adapter[netip.Addr] for IPv6 lookups.
type Adapter struct{}

// NewAdapter creates a new IPv6 adapter.
func NewAdapter() Adapter {
	return Adapter{}
}

// Encode converts an IPv6 address to a 128-bit binary string.
// IPv4-mapped IPv6 addresses (e.g. ::ffff:1.2.3.4) are rejected.
func (Adapter) Encode(addr netip.Addr) (string, error) {
	if !addr.Is6() || addr.Is4() || addr.Is4In6() {
		return "", fmt.Errorf("%w: %s", ErrNotIPv6, addr)
	}

	octets := addr.As16()

	var str []byte
	for _, octet := range octets {
		str = fmt.Appendf(str, "%08b", octet)
	}

	return string(str), nil
}

// KeyBits returns 128 for IPv6.
func (Adapter) KeyBits() uint {
	return 128
}

// ParseCIDR converts an IPv6 CIDR string to a binary prefix string for fltrie.Add.
// IPv4-mapped IPv6 addresses are rejected.
func ParseCIDR(cidr string) (string, error) {
	prefix, err := netip.ParsePrefix(cidr)
	if err != nil {
		return "", fmt.Errorf("cidr is not valid: %w", err)
	}

	addr := prefix.Addr()
	if !addr.Is6() || addr.Is4() || addr.Is4In6() {
		return "", fmt.Errorf("%w: %s", ErrNotIPv6, cidr)
	}

	octets := addr.As16()

	var str []byte
	for _, octet := range octets {
		str = fmt.Appendf(str, "%08b", octet)
	}

	return string(str)[:prefix.Bits()] + "*", nil
}

// DefaultConfig returns the standard FlashTrie config for IPv6.
func DefaultConfig() fltrie.Config {
	return fltrie.Config{
		KeyBits:   128,
		Stride:    8,
		CompSize:  2,
		TrieDepth: 2,
	}
}
