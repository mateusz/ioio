package pathfinder

import (
	"container/list"
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/beefsack/go-astar"
	"github.com/mateusz/ioio/architecture"
	"github.com/mateusz/rtsian/piksele"
)

type Pathfinder struct {
	gw      *piksele.World
	nodeMap [][]*PathNode
	width   int
}

type PathNode struct {
	X    int
	Y    int
	Cost float64
	from pathVec
	to   pathVec
	pf   *Pathfinder
}

type pathVec struct {
	x int
	y int
}

// Special handling for start, so that the pathfinder can get out
type pathStartingNode struct {
	x  int
	y  int
	pf *Pathfinder
}

func NewPathfinder(w *piksele.World, cs []*architecture.Component) Pathfinder {
	pf := Pathfinder{
		gw:    w,
		width: w.Tiles.Width,
	}
	pf.buildMap(cs)
	return pf
}

func (pf *Pathfinder) FindPath(fx, fy, tx, ty int) (l *list.List) {

	start := &pathStartingNode{
		x:  fx,
		y:  fy,
		pf: pf,
	}

	var path []astar.Pather
	found := false
	// Run through all targets. This is redundant, but works for now.
	// Idea: add pathTerminatingNodes so that the search only has to run once.
	for _, tn := range pf.getPatherNodesAt(tx, ty) {
		path, _, found = astar.Path(start, tn)
		if found {
			break
		}
	}

	if !found {
		return
	}

	l = list.New()
	for _, n := range path {
		l.PushFront(n)
	}
	// Remove starting tile
	l.Remove(l.Front())

	return
}

func (pf *Pathfinder) buildMap(cs []*architecture.Component) {
	cmap := make([]*architecture.Component, pf.gw.Tiles.Width*pf.gw.Tiles.Height)
	pf.nodeMap = make([][]*PathNode, pf.gw.Tiles.Width*pf.gw.Tiles.Height)

	for _, c := range cs {
		i := pf.gw.Tiles.Width*c.Y + c.X
		cmap[i] = c
	}

	for y := 0; y < pf.gw.Tiles.Height; y++ {
		for x := 0; x < pf.gw.Tiles.Width; x++ {
			i := x + y*pf.gw.Tiles.Width
			c := cmap[i]
			if cmap[i] == nil {
				continue
			}

			pf.nodeMap[i] = make([]*PathNode, 0)

			pf.nodeMap[i] = pf.convertToLinkages(c)
		}
	}
}

func (pf *Pathfinder) convertToLinkages(c *architecture.Component) []*PathNode {
	ns := make([]*PathNode, 0)
	var latMs int
	fmt.Sscanf(c.Lat, "%dms", &latMs)

	linkages := strings.Split(c.Con, ",")
	for _, l := range linkages {
		if len(l) != 2 {
			continue
		}
		n := &PathNode{
			Cost: float64(latMs),
			from: letterToDir(l[0]),
			to:   letterToDir(l[1]),
			X:    c.X,
			Y:    c.Y,
			pf:   pf,
		}
		ns = append(ns, n)
	}

	return ns
}

func letterToDir(letter byte) pathVec {
	switch letter {
	case 'x':
		// Host nodes can be gotten in, but not out. This is to stop pass-through traffic.
		return pathVec{x: 0, y: 0}
	case 'l':
		return pathVec{x: -1, y: 0}
	case 'r':
		return pathVec{x: 1, y: 0}
	case 't':
		return pathVec{x: 0, y: -1}
	case 'b':
		return pathVec{x: 0, y: 1}
	default:
		fmt.Printf("Unrecognised direction: %b", letter)
		os.Exit(2)
	}

	return pathVec{}
}

func (pf *Pathfinder) getPatherNodesAt(x, y int) []*PathNode {
	return pf.nodeMap[y*pf.width+x]
}

func (n *PathNode) PathNeighbors() []astar.Pather {
	ns := []astar.Pather{}
	if n.to.x < 0 && n.X > 0 {
		tns := n.pf.getPatherNodesAt(n.X-1, n.Y)
		for _, tn := range tns {
			if tn.from.x > 0 {
				ns = append(ns, tn)
			}
		}
	}
	if n.to.x > 0 && n.X < n.pf.gw.Tiles.Width-1 {
		tns := n.pf.getPatherNodesAt(n.X+1, n.Y)
		for _, tn := range tns {
			if tn.from.x < 0 {
				ns = append(ns, tn)
			}
		}
	}
	if n.to.y < 0 && n.Y > 0 {
		tns := n.pf.getPatherNodesAt(n.X, n.Y-1)
		for _, tn := range tns {
			if tn.from.y > 0 {
				ns = append(ns, tn)
			}
		}
	}
	if n.to.y > 0 && n.Y < n.pf.gw.Tiles.Height-1 {
		tns := n.pf.getPatherNodesAt(n.X, n.Y+1)
		for _, tn := range tns {
			if tn.from.y < 0 {
				ns = append(ns, tn)
			}
		}
	}
	return ns
}

func (n *PathNode) PathNeighborCost(to astar.Pather) float64 {
	tn, ok := to.(*PathNode)
	if !ok {
		return 10000000.0
	}

	return tn.Cost
}

func (n *PathNode) PathEstimatedCost(to astar.Pather) float64 {
	tn, ok := to.(*PathNode)
	if !ok {
		return 10000000.0
	}

	return math.Abs(float64(tn.X-n.X)) + math.Abs(float64(tn.Y-n.Y))
}

func (n *pathStartingNode) PathNeighbors() []astar.Pather {
	ns := []astar.Pather{}
	if n.x > 0 {
		tns := n.pf.getPatherNodesAt(n.x-1, n.y)
		for _, tn := range tns {
			if tn.from.x > 0 {
				ns = append(ns, tn)
			}
		}
	}
	if n.x < n.pf.gw.Tiles.Width-1 {
		tns := n.pf.getPatherNodesAt(n.x+1, n.y)
		for _, tn := range tns {
			if tn.from.x < 0 {
				ns = append(ns, tn)
			}
		}
	}
	if n.y > 0 {
		tns := n.pf.getPatherNodesAt(n.x, n.y+1)
		for _, tn := range tns {
			if tn.from.y > 0 {
				ns = append(ns, tn)
			}
		}
	}
	if n.y < n.pf.gw.Tiles.Height-1 {
		tns := n.pf.getPatherNodesAt(n.x, n.y+1)
		for _, tn := range tns {
			if tn.from.y < 0 {
				ns = append(ns, tn)
			}
		}
	}
	return ns
}

func (n *pathStartingNode) PathNeighborCost(to astar.Pather) float64 {
	tn, ok := to.(*PathNode)
	if !ok {
		return 10000000.0
	}

	return tn.Cost
}

func (n *pathStartingNode) PathEstimatedCost(to astar.Pather) float64 {
	tn, ok := to.(*pathStartingNode)
	if !ok {
		return 10000000.0
	}

	return math.Abs(float64(tn.x-n.x)) + math.Abs(float64(tn.y-n.y))
}
