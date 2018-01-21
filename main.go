/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 09-11-2017
 * |
 * | File Name:     main.go
 * +===============================================
 */

package main

import (
	"fmt"
	"log"
	"net"
)

func parseCIDR(cidr string) string {
	_, ipv4Net, err := net.ParseCIDR(cidr)
	if err != nil {
		log.Fatal(err)
	}

	size, _ := ipv4Net.Mask.Size()
	str := ""

	for _, octet := range ipv4Net.IP.To4() {
		str = str + fmt.Sprintf("%08b", octet)
	}
	str = str[:size] + "*"

	return str
}

func parseIP(ip string) string {
	ipv4 := net.ParseIP(ip)

	str := ""

	for _, octet := range ipv4.To4() {
		str = str + fmt.Sprintf("%08b", octet)
	}

	return str
}

func main() {
}
