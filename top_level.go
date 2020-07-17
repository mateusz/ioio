package main

import (
	"log"

	"golang.org/x/image/colornames"
)

type topLevel struct {
	get
	name string
}

func (tl *topLevel) exec() {
	h := tl.get.ctl["h"]
	dest := tl.program.findHostByName(h)
	gameBlips.add(dest.component.x, dest.component.y, colornames.Red)
	if dest == nil {
		log.Printf("[%s] Host not found '%s'", tl.name, h)
		return
	}
	tl.get.prg.exec(*dest)
	gameBlips.del(dest.component.x, dest.component.y, colornames.Red)
}
