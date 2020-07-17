package main

import "log"

type topLevel struct {
	get
	name string
}

func (tl *topLevel) exec() {
	h := tl.get.ctl["h"]
	dest := tl.program.findHostByName(h)
	// TODO add blip at dest
	if dest == nil {
		log.Printf("[%s] Host not found '%s'", tl.name, h)
		return
	}
	tl.get.prg.exec(*dest)
	// TODO remove blip at dest
}
