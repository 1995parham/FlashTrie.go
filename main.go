package main

import (
	"fmt"
	"log"
	"os"

	"github.com/1995parham/FlashTrie.go/fltrie"
	"github.com/1995parham/FlashTrie.go/net"
	"github.com/abiosoft/ishell"
	yaml "gopkg.in/yaml.v2"
)

type route struct {
	Route   string
	Nexthop string
}

func main() {
	var f string

	var routes []route

	fmt.Print("Routes file name (without .yml extension): ")
	fmt.Scanf("%s", &f)
	f += ".yml"

	data, err := os.ReadFile(f)
	if err != nil {
		log.Fatalf("Reading file %s failed with: %s\n", f, err)
	}

	err = yaml.Unmarshal(data, &routes)
	if err != nil {
		log.Fatalf("Parsing file %s failed with: %s\n", f, err)
	}

	// Building flash trie

	fltrie := fltrie.New()

	for _, route := range routes {
		r, err := net.ParseNet(route.Route)
		if err != nil {
			log.Fatalf("Parsing file %s failed with: %s\n", f, err)
		}

		fltrie.Add(r, route.Nexthop)
	}

	if err := fltrie.Build(); err != nil {
		log.Fatalf("Building flash trie failed with: %s\n", err)
	}

	// Run shell
	shell := ishell.New()
	shell.Println("Welcome to FlashTrie.go developed by Parham Alvani @ 2018")

	shell.AddCmd(&ishell.Cmd{
		Name: "lookup",
		Help: "lookup destination in routing table",
		Func: func(c *ishell.Context) {
			for _, arg := range c.Args {
				ip := net.ParseIP(arg)
				c.Printf("%s -> %s\n", arg, fltrie.Lookup(ip))
			}
		},
	})

	shell.Run()
}
