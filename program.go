package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type program struct {
	source       map[interface{}]interface{}
	top          []topLevel
	roComponents []*component
}

type topLevel struct {
	get
	name string
}

type get struct {
	program *program
	ctl     ctl
	prg     prg
}

type ctl map[string]string

type prg []interface{}

type compute int

func newProgram(fileName string, cs []*component) program {
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error loading program: %s\n", err)
		os.Exit(2)
	}
	yd := yaml.NewDecoder(f)

	p := program{
		source: make(map[interface{}]interface{}),
	}
	err = yd.Decode(p.source)
	if err != nil {
		fmt.Printf("Error parsing program: %s\n", err)
		os.Exit(2)
	}

	copy(p.roComponents, cs)

	p.parse()
	return p
}

func (p program) checkErr(ok bool, msg string) {
	if !ok {
		fmt.Println(msg)
		os.Exit(2)
	}
}

func (p *program) parse() {
	for k, v := range p.source {
		ks, ok := k.(string)
		p.checkErr(ok, "Error parsing top-level, hash key not a string")

		get, ok := v.(map[interface{}]interface{})
		p.checkErr(ok, "Error parsing top-level, not a hash")

		tl := topLevel{
			get:  p.parseGet(get),
			name: ks,
		}
		p.top = append(p.top, tl)
	}
}

func (p *program) parseGet(rawGet map[interface{}]interface{}) get {
	get := get{}
	get.program = p
	for k, v := range rawGet {
		ks, ok := k.(string)
		p.checkErr(ok, "Error parsing get, hash key not a string")

		if ks == "ctl" {
			ctl, ok := v.(map[interface{}]interface{})
			p.checkErr(ok, "Error parsing get, ctl not a hash")
			get.ctl = p.parseCtl(ctl)
		} else if ks == "prg" {
			prg, ok := v.([]interface{})
			p.checkErr(ok, "Error parsing get, prg not a list")
			get.prg = p.parsePrg(prg)
		}
	}

	return get
}

func (p *program) parseCtl(rawCtl map[interface{}]interface{}) ctl {
	ctl := make(ctl)
	for k, v := range rawCtl {
		ks, ok := k.(string)
		p.checkErr(ok, "Error parsing ctl, key not a string")
		vs, ok := v.(string)
		p.checkErr(ok, "Error parsing ctl, value not a string")

		ctl[ks] = vs
	}
	return ctl
}

func (p *program) parsePrg(rawPrg []interface{}) prg {
	prg := make(prg, 0)
	for _, stmt := range rawPrg {
		switch v := stmt.(type) {
		case string:
			if strings.HasPrefix(v, "c/") {
				var cAmount int
				n, err := fmt.Sscanf(v, "c/%d", &cAmount)
				if n != 1 {
					fmt.Printf("Error parsing prg instruction, c/ encountered too many times")
					os.Exit(2)
				}
				if err != nil {
					fmt.Printf("Error parsing prg instruction: %s\n", err)
					os.Exit(2)
				}
				prg = append(prg, compute(cAmount))
			}
		}
	}

	return prg
}

func (p *program) start() {
	for _, tl := range p.top {
		log.Printf("Starting top-level block '%s'", tl.name)
		tl.get.exec()
	}
}

func (g *get) exec() {
	go func() {
		h := g.ctl["h"]
		g.program.findByName(h)
		// TODO
	}()
}

func (p *program) findByName(name string) *component {
	for _, c := range p.roComponents {
		if c.name == name {
			return c
		}
	}
	return nil
}
