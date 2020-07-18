package main

import (
	"log"
)

type instruction interface {
	exec(prgScheduler)
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

	b := &blip{x: host.component.x, y: host.component.y, color: p.topLevel.color}
	gameBlips.add(b)

	for _, instr := range p.instructions {
		switch i := instr.(type) {
		case compute:
			i.exec(host.scheduler)
		}
	}

	gameBlips.del(b)

	log.Printf("[%s] Program ended on '%s'", p.topLevel.name, host.component.name)
}

type compute struct {
	c int
}

func (c compute) exec(sched prgScheduler) {
	sched.schedule(c.c)
}
