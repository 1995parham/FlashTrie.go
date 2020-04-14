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
	"io/ioutil"
	"log"

	"github.com/AUTProjects/FlashTrie.go/fltrie"
	"github.com/AUTProjects/FlashTrie.go/util"
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

	data, err := ioutil.ReadFile(f)
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
		r, err := util.ParseNet(route.Route)
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
	shell.Println("Welcome to FlashTrie.go by Parham Alvani @ 2018")

	shell.AddCmd(&ishell.Cmd{
		Name: "lookup",
		Help: "lookup destination in routing table",
		Func: func(c *ishell.Context) {
			for _, arg := range c.Args {
				ip := util.ParseIP(arg)
				c.Printf("%s -> %s\n", arg, fltrie.Lookup(ip))
			}
		},
	})

	shell.Run()
}
