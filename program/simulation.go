package program

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mateusz/ioio/pathfinder"
	"github.com/mateusz/rtsian/piksele"

	"github.com/mateusz/ioio/architecture"
	"github.com/mateusz/ioio/graphics"
	"golang.org/x/image/colornames"
	"gopkg.in/yaml.v2"
)

type Simulation struct {
	source     map[interface{}]interface{}
	top        []topLevel
	hosts      []host
	blipList   *graphics.BlipList
	pathfinder pathfinder.Pathfinder
}

type host struct {
	component architecture.Component
	scheduler hostScheduler
}

type hostScheduler interface {
	Schedule(time.Duration)
}

func NewSimulation(fileName string, cs []*architecture.Component, gw *piksele.World, bl *graphics.BlipList) Simulation {
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error loading simulation: %s\n", err)
		os.Exit(2)
	}
	yd := yaml.NewDecoder(f)

	s := Simulation{
		source:     make(map[interface{}]interface{}),
		blipList:   bl,
		pathfinder: pathfinder.NewPathfinder(gw, cs),
	}
	err = yd.Decode(s.source)
	if err != nil {
		fmt.Printf("Error parsing simulation: %s\n", err)
		os.Exit(2)
	}

	s.hosts = make([]host, 0)
	s.copyComponents(cs)

	s.parse()
	return s
}

func (s *Simulation) copyComponents(cs []*architecture.Component) {
	for _, c := range cs {
		h := host{
			component: *c,
		}

		switch c.Sched {
		case "infinite":
			h.scheduler = architecture.NewSchedInfinite(&h.component)
		case "multitasking":
			h.scheduler = architecture.NewSchedMultitasking(&h.component)
		}

		s.hosts = append(s.hosts, h)
	}
}

func (s *Simulation) checkErr(ok bool, msg string) {
	if !ok {
		fmt.Println(msg)
		os.Exit(2)
	}
}

func (s *Simulation) parse() {
	for k, v := range s.source {
		ks, ok := k.(string)
		s.checkErr(ok, "Error parsing top-level, hash key not a string")

		get, ok := v.(map[interface{}]interface{})
		s.checkErr(ok, "Error parsing top-level, not a hash")

		tl := topLevel{
			name: ks,
		}
		// Top-level is a kind of a get.
		tl.get = s.parseGet(&tl, get)
		if colorName, ok := tl.get.ctl["color"]; ok {
			if tl.color, ok = colornames.Map[colorName]; !ok {
				s.checkErr(ok, fmt.Sprintf("Unrecognised color name: '%s'", colorName))
			}
		}

		s.top = append(s.top, tl)
	}
}

func (s *Simulation) parseGet(tl *topLevel, rawGet map[interface{}]interface{}) get {
	get := get{
		sim:      s,
		topLevel: tl,
	}
	for k, v := range rawGet {
		ks, ok := k.(string)
		s.checkErr(ok, "Error parsing get, hash key not a string")

		if ks == "ctl" {
			ctl, ok := v.(map[interface{}]interface{})
			s.checkErr(ok, "Error parsing get, ctl not a hash")
			get.ctl = s.parseCtl(ctl)
		} else if ks == "prg" {
			prg, ok := v.([]interface{})
			s.checkErr(ok, "Error parsing get, prg not a list")
			get.prg = s.parsePrg(tl, prg)
		}
	}

	return get
}

func (s *Simulation) parseSerial(tl *topLevel, rawSerial map[interface{}]interface{}) serial {
	serial := serial{
		sim:      s,
		topLevel: tl,
	}
	for k, v := range rawSerial {
		ks, ok := k.(string)
		s.checkErr(ok, "Error parsing serial, hash key not a string")

		if ks == "ctl" {
			ctl, ok := v.(map[interface{}]interface{})
			s.checkErr(ok, "Error parsing serial, ctl not a hash")
			serial.ctl = s.parseCtl(ctl)
		} else if ks == "prg" {
			prg, ok := v.([]interface{})
			s.checkErr(ok, "Error parsing serial, prg not a list")
			serial.prg = s.parsePrg(tl, prg)
		}
	}

	return serial
}

func (s *Simulation) parseParallel(tl *topLevel, rawParallel map[interface{}]interface{}) parallel {
	parallel := parallel{
		sim:      s,
		topLevel: tl,
	}
	for k, v := range rawParallel {
		ks, ok := k.(string)
		s.checkErr(ok, "Error parsing parallel, hash key not a string")

		if ks == "ctl" {
			ctl, ok := v.(map[interface{}]interface{})
			s.checkErr(ok, "Error parsing parallel, ctl not a hash")
			parallel.ctl = s.parseCtl(ctl)
		} else if ks == "prg" {
			prg, ok := v.([]interface{})
			s.checkErr(ok, "Error parsing parallel, prg not a list")
			parallel.prg = s.parsePrg(tl, prg)
		}
	}

	return parallel
}

func (s *Simulation) parseRps(tl *topLevel, rawRps map[interface{}]interface{}) rps {
	rps := rps{
		sim:      s,
		topLevel: tl,
	}
	for k, v := range rawRps {
		ks, ok := k.(string)
		s.checkErr(ok, "Error parsing rps, hash key not a string")

		if ks == "ctl" {
			ctl, ok := v.(map[interface{}]interface{})
			s.checkErr(ok, "Error parsing rps, ctl not a hash")
			rps.ctl = s.parseCtl(ctl)
		} else if ks == "prg" {
			prg, ok := v.([]interface{})
			s.checkErr(ok, "Error parsing rps, prg not a list")
			rps.prg = s.parsePrg(tl, prg)
		}
	}

	return rps
}

func (s *Simulation) parseCtl(rawCtl map[interface{}]interface{}) ctl {
	ctl := make(ctl)
	for k, v := range rawCtl {
		ks, ok := k.(string)
		s.checkErr(ok, "Error parsing ctl, key not a string")
		vs, ok := v.(string)
		s.checkErr(ok, "Error parsing ctl, value not a string")

		ctl[ks] = vs
	}
	return ctl
}

func (s *Simulation) parsePrg(tl *topLevel, rawPrg []interface{}) prg {
	prg := prg{
		instructions: make([]instruction, 0, len(rawPrg)),
		topLevel:     tl,
	}
	for _, stmt := range rawPrg {
		switch v := stmt.(type) {
		case string:
			if strings.HasPrefix(v, "c/") {
				var cAmount string
				n, err := fmt.Sscanf(v, "c/%s", &cAmount)
				if n != 1 {
					fmt.Printf("Error parsing prg instruction, c/ encountered too many times")
					os.Exit(2)
				}
				if err != nil {
					fmt.Printf("Error parsing prg instruction: %s\n", err)
					os.Exit(2)
				}

				c, err := time.ParseDuration(cAmount)
				if err != nil {
					fmt.Printf("Error parsing prg instruction: %s\n", err)
					os.Exit(2)
				}
				prg.instructions = append(prg.instructions, compute{
					topLevel: tl,
					c:        c,
				})
			}
		case map[interface{}]interface{}:
			for mk, mv := range v {
				mks, ok := mk.(string)
				s.checkErr(ok, "Error parsing get, instruction not a hash")

				if strings.HasPrefix(mks, "get") {
					gv, ok := mv.(map[interface{}]interface{})
					s.checkErr(ok, "Error parsing get, not a hash")

					prg.instructions = append(prg.instructions, s.parseGet(tl, gv))
				} else if strings.HasPrefix(mks, "serial") {
					gv, ok := mv.(map[interface{}]interface{})
					s.checkErr(ok, "Error parsing serial, not a hash")

					prg.instructions = append(prg.instructions, s.parseSerial(tl, gv))
				} else if strings.HasPrefix(mks, "parallel") {
					gv, ok := mv.(map[interface{}]interface{})
					s.checkErr(ok, "Error parsing parallel, not a hash")

					prg.instructions = append(prg.instructions, s.parseParallel(tl, gv))
				} else if strings.HasPrefix(mks, "rps") {
					gv, ok := mv.(map[interface{}]interface{})
					s.checkErr(ok, "Error parsing rps, not a hash")

					prg.instructions = append(prg.instructions, s.parseRps(tl, gv))
				}
			}
		}
	}

	return prg
}

func (s *Simulation) Start() {
	for _, tl := range s.top {
		tlLocal := tl
		go tlLocal.exec()
	}
}

func (s *Simulation) findHostByName(name string) *host {
	for _, h := range s.hosts {
		if h.component.Name == name {
			return &h
		}
	}
	return nil
}
