package main

import (
	"log"
)

type program struct {
	source    map[interface{}]interface{}
	top       []topLevel
	hosts     []host
	presences []presentable
}

type host struct {
	component component
	scheduler prgScheduler
}

type prgScheduler interface {
	schedule(int)
}

type ctl map[string]string

type instruction interface {
	exec(prgScheduler)
}

type compute struct {
	c int
}

func (p *program) start() {
	for _, tl := range p.top {
		log.Printf("[%s] Starting request", tl.name)
		go tl.exec()
	}
}

func (c compute) exec(sched prgScheduler) {
	sched.schedule(c.c)
}

func (p *program) findHostByName(name string) *host {
	for _, h := range p.hosts {
		if h.component.name == name {
			return &h
		}
	}
	return nil
}
