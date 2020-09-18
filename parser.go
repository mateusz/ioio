package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/image/colornames"

	"gopkg.in/yaml.v2"
)

type ctl map[string]string

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

		switch c.sched {
		case "infinite":
			h.scheduler = NewSchedInfinite(&h.component)
		case "multitasking":
			h.scheduler = NewSchedMultitasking(&h.component)
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

func (p *program) parseSerial(tl *topLevel, rawSerial map[interface{}]interface{}) serial {
	serial := serial{
		program:  p,
		topLevel: tl,
	}
	for k, v := range rawSerial {
		ks, ok := k.(string)
		p.checkErr(ok, "Error parsing serial, hash key not a string")

		if ks == "ctl" {
			ctl, ok := v.(map[interface{}]interface{})
			p.checkErr(ok, "Error parsing serial, ctl not a hash")
			serial.ctl = p.parseCtl(ctl)
		} else if ks == "prg" {
			prg, ok := v.([]interface{})
			p.checkErr(ok, "Error parsing serial, prg not a list")
			serial.prg = p.parsePrg(tl, prg)
		}
	}

	return serial
}

func (p *program) parseParallel(tl *topLevel, rawParallel map[interface{}]interface{}) parallel {
	parallel := parallel{
		program:  p,
		topLevel: tl,
	}
	for k, v := range rawParallel {
		ks, ok := k.(string)
		p.checkErr(ok, "Error parsing parallel, hash key not a string")

		if ks == "ctl" {
			ctl, ok := v.(map[interface{}]interface{})
			p.checkErr(ok, "Error parsing parallel, ctl not a hash")
			parallel.ctl = p.parseCtl(ctl)
		} else if ks == "prg" {
			prg, ok := v.([]interface{})
			p.checkErr(ok, "Error parsing parallel, prg not a list")
			parallel.prg = p.parsePrg(tl, prg)
		}
	}

	return parallel
}

func (p *program) parseRps(tl *topLevel, rawRps map[interface{}]interface{}) rps {
	rps := rps{
		program:  p,
		topLevel: tl,
	}
	for k, v := range rawRps {
		ks, ok := k.(string)
		p.checkErr(ok, "Error parsing rps, hash key not a string")

		if ks == "ctl" {
			ctl, ok := v.(map[interface{}]interface{})
			p.checkErr(ok, "Error parsing rps, ctl not a hash")
			rps.ctl = p.parseCtl(ctl)
		} else if ks == "prg" {
			prg, ok := v.([]interface{})
			p.checkErr(ok, "Error parsing rps, prg not a list")
			rps.prg = p.parsePrg(tl, prg)
		}
	}

	return rps
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
				prg.instructions = append(prg.instructions, compute{
					topLevel: tl,
					c:        cAmount,
				})
			}
		case map[interface{}]interface{}:
			for mk, mv := range v {
				mks, ok := mk.(string)
				p.checkErr(ok, "Error parsing get, instruction not a hash")

				if strings.HasPrefix(mks, "get") {
					gv, ok := mv.(map[interface{}]interface{})
					p.checkErr(ok, "Error parsing get, not a hash")

					prg.instructions = append(prg.instructions, p.parseGet(tl, gv))
				} else if strings.HasPrefix(mks, "serial") {
					gv, ok := mv.(map[interface{}]interface{})
					p.checkErr(ok, "Error parsing serial, not a hash")

					prg.instructions = append(prg.instructions, p.parseSerial(tl, gv))
				} else if strings.HasPrefix(mks, "parallel") {
					gv, ok := mv.(map[interface{}]interface{})
					p.checkErr(ok, "Error parsing parallel, not a hash")

					prg.instructions = append(prg.instructions, p.parseParallel(tl, gv))
				} else if strings.HasPrefix(mks, "rps") {
					gv, ok := mv.(map[interface{}]interface{})
					p.checkErr(ok, "Error parsing rps, not a hash")

					prg.instructions = append(prg.instructions, p.parseRps(tl, gv))
				}
			}
		}
	}

	return prg
}
