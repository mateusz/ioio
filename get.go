package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"time"
)

type get struct {
	program  *program
	topLevel *topLevel
	ctl      ctl
	prg      prg
	color    *color.Color
}

func (g get) exec(origin host) {
	h := g.ctl["h"]
	dest := g.program.findHostByName(h)
	if dest == nil {
		fmt.Printf("Host '%s' not found\n", h)
		os.Exit(2)
	}

	log.Printf("[%s] Starting get '%s'\n", g.topLevel.name, dest.component.name)
	b := &blip{x: origin.component.x, y: origin.component.y, color: g.topLevel.color}
	gameBlips.add(b)

	if !g.transit(origin, *dest) {
		return
	}

	g.prg.exec(*dest)

	// But what if we can't find the return path? Should we instead simulate
	// stateful firewalls, and make the return path simply the reverse?
	// But then what if a wire breaks on the way back? That wouldn't be simulatable.
	g.transit(*dest, origin)

	gameBlips.del(b)
	log.Printf("[%s] Finished get '%s'\n", g.topLevel.name, dest.component.name)
}

func (g get) transit(from host, to host) bool {
	path := gamePathfinder.findPath(
		pathVec{x: from.component.x, y: from.component.y},
		pathVec{x: to.component.x, y: to.component.y},
	)
	if path == nil || path.Len() == 0 {
		log.Printf("[%s] Path not found from '%s' to '%s'\n", g.topLevel.name, from.component.name, to.component.name)
		return false
	}

	// TODO move inside findPath
	var slat int
	fmt.Sscanf(from.component.lat, "%dms", &slat)
	totalCost := 0.0 + float64(slat)
	for e := path.Front(); e != nil; e = e.Next() {
		pn, ok := e.Value.(*pathNode)
		if !ok {
			fmt.Print("Non-pathNode found in path list\n")
			os.Exit(2)
		}

		totalCost += pn.cost
	}

	b := &blip{
		x:     from.component.x,
		y:     from.component.y,
		color: g.topLevel.color,
		path:  path, // takes ownership
	}
	gameBlips.add(b)

	// Now traveling for as long as it takes, blip will take care of the actual animation
	// Offsetting by 50 prevents the final blink that makes the dot jump to "from". I haven't
	// debugged this one yet, and it seems to happen more on the return path.
	time.Sleep(time.Duration(totalCost-50.0) * time.Millisecond)

	gameBlips.del(b)
	return true
}
