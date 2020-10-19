package program

import (
	"log"
	"sync"
)

type parallel struct {
	sim      *Simulation
	topLevel *topLevel
	ctl      ctl
	prg      prg
}

func (p parallel) exec(origin host) {
	log.Printf("[%s] Starting parallel execution\n", p.topLevel.name)

	var wg sync.WaitGroup
	r := p.ctl.int("r")
	wg.Add(r)
	for i := 0; i < r; i++ {
		go func(p parallel, origin host, wg *sync.WaitGroup) {
			p.prg.exec(origin)
			wg.Done()
		}(p, origin, &wg)
	}
	wg.Wait()

	log.Printf("[%s] Finished parallel execution\n", p.topLevel.name)
}
