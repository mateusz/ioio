package main

import "log"

type prg struct {
	program      *program
	topLevel     *topLevel
	instructions []instruction
}

func (p *prg) exec(host host) {
	log.Printf("[%s] Running program on '%s'", p.topLevel.name, host.component.name)
	for _, instr := range p.instructions {
		switch i := instr.(type) {
		case compute:
			i.exec(host.scheduler)
		}
	}
	log.Printf("[%s] Program ended on '%s'", p.topLevel.name, host.component.name)
}
