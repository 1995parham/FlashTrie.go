/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 22-01-2018
 * |
 * | File Name:     util/util.go
 * +===============================================
 */

package util

import (
	"fmt"
	"net"
)

// ParseNet parse given v4 network into binary string.
// 192.0.0.0/4 -> 1100*
func ParseNet(cidr string) (string, error) {
	_, ipv4Net, err := net.ParseCIDR(cidr)
	if err != nil {
		return "", err
	}

	size, _ := ipv4Net.Mask.Size()
	str := ""

	for _, octet := range ipv4Net.IP.To4() {
		str = str + fmt.Sprintf("%08b", octet)
	}
	str = str[:size] + "*"

	return str, nil
}

// ParseIP parse given ip v4 into binary string.
func ParseIP(ip string) string {
	ipv4 := net.ParseIP(ip)

	str := ""

	for _, octet := range ipv4.To4() {
		str = str + fmt.Sprintf("%08b", octet)
	}

	return str
}
