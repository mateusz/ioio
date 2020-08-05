package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
)

type parallel struct {
	program  *program
	topLevel *topLevel
	ctl      ctl
	prg      prg
}

func (p parallel) exec(origin host) {
	log.Printf("[%s] Starting parallel execution\n", p.topLevel.name)
	r, err := strconv.Atoi(p.ctl["r"])
	if err != nil {
		fmt.Printf("[%s] Failed to convert ctl.r: %s\n", p.topLevel.name, err)
		os.Exit(2)
	}

	var wg sync.WaitGroup
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
