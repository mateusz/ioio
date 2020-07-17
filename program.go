package main

import (
	"log"
)

type program struct {
	source map[interface{}]interface{}
	top    []topLevel
	hosts  []host
}

type host struct {
	component component
	scheduler hostScheduler
}

type hostScheduler interface {
	schedule(int)
}

func (p *program) start() {
	for _, tl := range p.top {
		log.Printf("[%s] Starting request", tl.name)
		go tl.exec()
	}
}

func (p *program) findHostByName(name string) *host {
	for _, h := range p.hosts {
		if h.component.name == name {
			return &h
		}
	}
	return nil
}
