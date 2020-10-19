package program

import (
	"fmt"
	"log"
	"os"
	"strconv"
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
	ctlR, err := strconv.ParseFloat(r.ctl["r"], 64)
	if err != nil {
		fmt.Printf("[%s] Failed to convert ctl.r: %s\n", r.topLevel.name, err)
		os.Exit(2)
	}

	interval := time.Duration(1.0 / ctlR * float64(time.Second))
	log.Printf("[%s] Interval set to %.2fs\n", r.topLevel.name, float64(interval)/float64(time.Second))
	for range time.Tick(interval) {
		go func(r rps, origin host) {
			r.prg.exec(origin)
		}(r, origin)
	}

	// Will never get here
	log.Printf("[%s] Finished rps execution\n", r.topLevel.name)
}
