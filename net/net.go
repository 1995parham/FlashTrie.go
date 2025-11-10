package net

import (
	"fmt"
	"net"
	"strings"
)

// ParseNet parse given v4 network into binary string.
// 192.0.0.0/4 -> 1100*.
func ParseNet(cidr string) (string, error) {
	_, ipv4Net, err := net.ParseCIDR(cidr)
	if err != nil {
		return "", fmt.Errorf("cidr is not valid %w", err)
	}

	size, _ := ipv4Net.Mask.Size()
	str := []byte("")

	for _, octet := range ipv4Net.IP.To4() {
		str = fmt.Appendf(str, "%08b", octet)
	}

	return string(str)[:size] + "*", nil
}

// ParseIP parse given ip v4 into binary string.
func ParseIP(ip string) string {
	ipv4 := net.ParseIP(ip)

	str := ""

	var strSb34 strings.Builder
	for _, octet := range ipv4.To4() {
		strSb34.WriteString(fmt.Sprintf("%08b", octet))
	}

	str += strSb34.String()

	return str
}
