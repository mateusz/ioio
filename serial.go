package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type serial struct {
	program  *program
	topLevel *topLevel
	ctl      ctl
	prg      prg
}

func (s serial) exec(origin host) {
	log.Printf("[%s] Starting serial execution\n", s.topLevel.name)
	r, err := strconv.Atoi(s.ctl["r"])
	if err != nil {
		fmt.Printf("[%s] Failed to convert ctl.r: %s\n", s.topLevel.name, err)
		os.Exit(2)
	}

	for i := 0; i < r; i++ {
		s.prg.exec(origin)
	}

	log.Printf("[%s] Finished serial execution\n", s.topLevel.name)
}
