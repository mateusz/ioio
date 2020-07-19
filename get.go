package main

import (
	"fmt"
	"image/color"
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

type ctl map[string]string

func (g get) exec(origin host) {
	h := g.ctl["h"]
	dest := g.program.findHostByName(h)
	if dest == nil {
		fmt.Printf("Host '%s' not found", h)
		os.Exit(2)
	}

	path := gamePathfinder.findPath(
		pathVec{x: origin.component.x, y: origin.component.y},
		pathVec{x: dest.component.x, y: dest.component.y},
	)
	if path == nil {
		return
	}

	totalCost := 0.0
	for e := path.Front(); e != nil; e = e.Next() {
		pn, ok := e.Value.(*pathNode)
		if !ok {
			fmt.Print("Non-pathNode found in path list\n")
			os.Exit(2)
		}

		totalCost += pn.cost
	}

	b := &blip{
		x:     origin.component.x,
		y:     origin.component.y,
		color: g.topLevel.color,
		path:  path,
	}
	gameBlips.add(b)

	// Now traveling for as long as it takes, blip will take care of the actual animation
	time.Sleep(time.Duration(totalCost) * time.Millisecond)

	gameBlips.del(b)

	g.prg.exec(*dest)
}
