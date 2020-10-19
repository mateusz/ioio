package program

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"time"

	"github.com/mateusz/ioio/graphics"
)

type pathNode interface {
	Cost() time.Duration
}

type get struct {
	sim      *Simulation
	topLevel *topLevel
	ctl      ctl
	prg      prg
	color    *color.Color
}

func (g get) exec(origin host) {
	h := g.ctl.string("h")
	dest := g.sim.findHostByName(h)
	if dest == nil {
		fmt.Printf("Host '%s' not found\n", h)
		os.Exit(2)
	}

	log.Printf("[%s] Starting get '%s'\n", g.topLevel.name, dest.component.Name)
	b := &graphics.Blip{X: origin.component.X, Y: origin.component.Y, Color: g.topLevel.color}
	g.sim.blipList.Give(b)
	defer g.sim.blipList.Del(b)

	if !g.transit(origin, *dest) {
		return
	}

	// Execute at dest
	g.prg.exec(*dest)

	// But what if we can't find the return path? Should we instead simulate
	// stateful firewalls, and make the return path simply the reverse?
	// But then what if a wire breaks on the way back? That wouldn't be simulatable.
	g.transit(*dest, origin)

	log.Printf("[%s] Finished get '%s'\n", g.topLevel.name, dest.component.Name)
}

func (g get) transit(from host, to host) bool {
	path := g.sim.pathfinder.FindPath(from.component.X, from.component.Y, to.component.X, to.component.Y)
	if path == nil || path.Len() == 0 {
		log.Printf("[%s] Path not found from '%s' to '%s'\n", g.topLevel.name, from.component.Name, to.component.Name)
		return false
	}

	var totalCost time.Duration
	for e := path.Front(); e != nil; e = e.Next() {
		pn, ok := e.Value.(pathNode)
		if !ok {
			fmt.Print("Non-pathNode found in path list (2)\n")
			os.Exit(2)
		}

		// Cost of origin tile ignored.
		if e == path.Front() {
			continue
		}
		totalCost += pn.Cost()
	}

	b := &graphics.Blip{
		X:     from.component.X,
		Y:     from.component.Y,
		Color: g.topLevel.color,
		Path:  path,
	}
	g.sim.blipList.Give(b)

	// The blip is being animated in a goroutine.
	time.Sleep(totalCost)

	g.sim.blipList.Del(b)
	return true
}
