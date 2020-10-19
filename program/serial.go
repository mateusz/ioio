package program

import (
	"log"
)

type serial struct {
	sim      *Simulation
	topLevel *topLevel
	ctl      ctl
	prg      prg
}

func (s serial) exec(origin host) {
	log.Printf("[%s] Starting serial execution\n", s.topLevel.name)

	r := s.ctl.int("r")
	for i := 0; i < r; i++ {
		s.prg.exec(origin)
	}

	log.Printf("[%s] Finished serial execution\n", s.topLevel.name)
}
