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

func main() {
	_, ipv4Net, err := net.ParseCIDR("192.0.2.1/24")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ipv4Net.IP[0])
	fmt.Println(ipv4Net.IP[1])
	fmt.Println(ipv4Net.IP[2])
	fmt.Println(ipv4Net.IP[3])
}
