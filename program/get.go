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
	Cost() float64
}

type get struct {
	sim      *Simulation
	topLevel *topLevel
	ctl      ctl
	prg      prg
	color    *color.Color
}

func (g get) exec(origin host) {
	h := g.ctl["h"]
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

	totalCost := 0.0
	for e := path.Front(); e != nil; e = e.Next() {
		pn, ok := e.Value.(pathNode)
		if !ok {
			fmt.Print("Non-pathNode found in path list\n")
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

	// Now traveling for as long as it takes, blip will take care of the actual animation
	// Offsetting by 50 prevents the final blink that makes the dot jump to "from". I haven't
	// debugged this one yet, and it seems to happen more on the return path.
	time.Sleep(time.Duration(totalCost) * time.Millisecond)

	g.sim.blipList.Del(b)
	return true
}
