package program

import (
	"image/color"
	"log"

	"github.com/mateusz/ioio/graphics"
)

type topLevel struct {
	get
	name  string
	color color.Color
}

func (tl *topLevel) exec() {
	h := tl.get.ctl["h"]
	dest := tl.sim.findHostByName(h)
	b := &graphics.Blip{X: dest.component.X, Y: dest.component.Y, Color: tl.color}
	tl.sim.blipList.Give(b)
	if dest == nil {
		log.Printf("[%s] Host not found '%s'", tl.name, h)
		return
	}
	tl.sim.blipList.Del(b)
	tl.get.prg.exec(*dest)
}
