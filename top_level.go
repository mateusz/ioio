package main

import "log"

type topLevel struct {
	get
	name string
}

func (tl *topLevel) exec() {
	h := tl.get.ctl["h"]
	dest := tl.program.findHostByName(h)
	if dest == nil {
		log.Printf("[%s] Host not found '%s'", tl.name, h)
		return
	}
	tl.get.prg.exec(*dest)
}
