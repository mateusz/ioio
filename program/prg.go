package program

import (
	"log"

	"github.com/mateusz/ioio/graphics"
)

type instruction interface {
	exec(host)
}

type prg struct {
	sim          *Simulation
	topLevel     *topLevel
	instructions []instruction
}

func (p *prg) exec(host host) {
	log.Printf("[%s] Running prg on '%s'", p.topLevel.name, host.component.Name)

	for _, instr := range p.instructions {
		instr.exec(host)
	}

	log.Printf("[%s] Prg ended on '%s'", p.topLevel.name, host.component.Name)
}

type compute struct {
	topLevel *topLevel
	c        int
}

func (c compute) exec(h host) {
	b := &graphics.Blip{X: h.component.X, Y: h.component.Y, Color: c.topLevel.color}
	c.topLevel.sim.blipList.Give(b)
	h.scheduler.Schedule(c.c)
	c.topLevel.sim.blipList.Del(b)
}
