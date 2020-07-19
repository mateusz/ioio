package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/image/colornames"

	"gopkg.in/yaml.v2"
)

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

	p.hosts = make([]host, 0)
	p.copyComponents(cs)

	p.parse()
	return p
}

func (p *program) copyComponents(cs []*component) {
	for _, c := range cs {
		h := host{
			component: *c,
		}

		if c.sched == "infinite" {
			h.scheduler = NewSchedInfinite(&h.component)
		}

		p.hosts = append(p.hosts, h)
	}
}

func (p *program) checkErr(ok bool, msg string) {
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
			name: ks,
		}
		// Top-level is a kind of a get.
		tl.get = p.parseGet(&tl, get)
		if colorName, ok := tl.get.ctl["color"]; ok {
			if tl.color, ok = colornames.Map[colorName]; !ok {
				p.checkErr(ok, fmt.Sprintf("Unrecognised color name: '%s'", colorName))
			}
		}

		p.top = append(p.top, tl)
	}
}

func (p *program) parseGet(tl *topLevel, rawGet map[interface{}]interface{}) get {
	get := get{
		program:  p,
		topLevel: tl,
	}
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
			get.prg = p.parsePrg(tl, prg)
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

func (p *program) parsePrg(tl *topLevel, rawPrg []interface{}) prg {
	prg := prg{
		instructions: make([]instruction, 0, len(rawPrg)),
		topLevel:     tl,
	}
	for _, stmt := range rawPrg {
		switch v := stmt.(type) {
		case string:
			if strings.HasPrefix(v, "c/") {
				var cAmount int
				n, err := fmt.Sscanf(v, "c/%dms", &cAmount)
				if n != 1 {
					fmt.Printf("Error parsing prg instruction, c/ encountered too many times")
					os.Exit(2)
				}
				if err != nil {
					fmt.Printf("Error parsing prg instruction: %s\n", err)
					os.Exit(2)
				}
				prg.instructions = append(prg.instructions, compute{c: cAmount})
			}
		case map[interface{}]interface{}:
			for mk, mv := range v {
				mks, ok := mk.(string)
				p.checkErr(ok, "Error parsing get, instruction not a hash")

				if strings.HasPrefix(mks, "get") {
					gv, ok := mv.(map[interface{}]interface{})
					p.checkErr(ok, "Error parsing get, get not a hash")

					prg.instructions = append(prg.instructions, p.parseGet(tl, gv))
				}
			}
		}
	}

	return prg
}
