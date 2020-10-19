package program

import (
	"log"
	"time"
)

type rps struct {
	sim      *Simulation
	topLevel *topLevel
	ctl      ctl
	prg      prg
}

func (r rps) exec(origin host) {
	log.Printf("[%s] Starting rps execution\n", r.topLevel.name)

	interval := time.Duration(1.0 / r.ctl.float64("r") * float64(time.Second))
	log.Printf("[%s] Interval set to %.2fs\n", r.topLevel.name, float64(interval)/float64(time.Second))
	for range time.Tick(interval) {
		go func(r rps, origin host) {
			r.prg.exec(origin)
		}(r, origin)
	}

	// Will never get here
	log.Printf("[%s] Finished rps execution\n", r.topLevel.name)
}
