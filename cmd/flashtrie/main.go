package main

import (
	"bufio"
	"fmt"
	"log"
	"net/netip"
	"os"
	"strings"

	"github.com/1995parham/FlashTrie.go/fltrie"
	"github.com/1995parham/FlashTrie.go/ipv4"
	yaml "gopkg.in/yaml.v2"
)

type route struct {
	Route   string
	Nexthop string
}

func main() {
	var f string

	var routes []route

	// nolint: forbidigo
	fmt.Print("Routes file name (without .yml extension): ")

	if n, err := fmt.Scanf("%s", &f); n != 1 || err != nil {
		log.Fatalf("invalid input")
	}

	f += ".yml"

	data, err := os.ReadFile(f)
	if err != nil {
		log.Fatalf("Reading file %s failed with: %s\n", f, err)
	}

	if err := yaml.Unmarshal(data, &routes); err != nil {
		log.Fatalf("Parsing file %s failed with: %s\n", f, err)
	}

	// Building flash trie
	fl := fltrie.New[netip.Addr, string](ipv4.NewAdapter(), ipv4.DefaultConfig())

	for _, route := range routes {
		r, err := ipv4.ParseCIDR(route.Route)
		if err != nil {
			log.Fatalf("Parsing file %s failed with: %s\n", f, err)
		}

		if err := fl.Add(r, route.Nexthop); err != nil {
			log.Fatalf("Adding route failed: %s\n", err)
		}
	}

	if err := fl.Build(); err != nil {
		log.Fatalf("Building flash trie failed with: %s\n", err)
	}

	repl(fl)
}

func repl(fl *fltrie.FLTrie[netip.Addr, string]) {
	// nolint: forbidigo
	fmt.Println("Welcome to FlashTrie.go developed by Parham Alvani @ 2018")
	// nolint: forbidigo
	fmt.Println("Commands: lookup <ip>..., exit")

	scanner := bufio.NewScanner(os.Stdin)

	// nolint: forbidigo
	fmt.Print(">>> ")

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())

		if len(fields) == 0 {
			// nolint: forbidigo
			fmt.Print(">>> ")

			continue
		}

		switch fields[0] {
		case "lookup":
			for _, arg := range fields[1:] {
				addr, err := netip.ParseAddr(arg)
				if err != nil {
					log.Printf("invalid IP %s: %s", arg, err)

					continue
				}

				result, found, err := fl.Lookup(addr)
				if err != nil {
					log.Printf("lookup error for %s: %s", arg, err)

					continue
				}

				if found {
					// nolint: forbidigo
					fmt.Printf("%s -> %s\n", arg, result)
				} else {
					// nolint: forbidigo
					fmt.Printf("%s -> (no match)\n", arg)
				}
			}
		case "exit", "quit":
			return
		default:
			// nolint: forbidigo
			fmt.Printf("unknown command: %s\n", fields[0])
		}

		// nolint: forbidigo
		fmt.Print(">>> ")
	}
}
