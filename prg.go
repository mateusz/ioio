package main

import (
	"log"
)

type instruction interface {
	exec(host)
}

type prgScheduler interface {
	schedule(int)
}

type prg struct {
	program      *program
	topLevel     *topLevel
	instructions []instruction
}

func (p *prg) exec(host host) {
	log.Printf("[%s] Running program on '%s'", p.topLevel.name, host.component.name)

	for _, instr := range p.instructions {
		instr.exec(host)
	}

	log.Printf("[%s] Program ended on '%s'", p.topLevel.name, host.component.name)
}

type compute struct {
	topLevel *topLevel
	c        int
}

func (c compute) exec(h host) {
	b := &blip{x: h.component.x, y: h.component.y, color: c.topLevel.color}
	gameBlips.add(b)
	h.scheduler.schedule(c.c)
	gameBlips.del(b)
}
