package main

import (
	"image/color"
	"log"
)

type topLevel struct {
	get
	name  string
	color color.Color
}

func (tl *topLevel) exec() {
	h := tl.get.ctl["h"]
	dest := tl.program.findHostByName(h)
	b := &blip{x: dest.component.x, y: dest.component.y, color: tl.color}
	gameBlips.add(b)
	if dest == nil {
		log.Printf("[%s] Host not found '%s'", tl.name, h)
		return
	}
	tl.get.prg.exec(*dest)
	gameBlips.del(b)
}
