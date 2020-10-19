package graphics

import (
	"container/list"
	"fmt"
	"image/color"
	"math"
	"os"
	"time"

	"github.com/faiface/pixel"
	"github.com/mateusz/rtsian/piksele"
)

type pathNode interface {
	X() int
	Y() int
	Cost() float64
}

type Blip struct {
	Color     color.Color
	X         int        // logical location
	Y         int        // logcial location
	Path      *list.List // current movement path to target
	Size      pixel.Vec  // pain size
	Pos       pixel.Vec  // paint location
	animStart time.Time  // start of current movement target
}

type BlipList struct {
	blips   list.List
	addChan chan *Blip
	delChan chan *Blip
	getChan chan chan []Blip
	gw      *piksele.World
}

func NewBlipList(gw *piksele.World) BlipList {
	bl := BlipList{
		addChan: make(chan *Blip),
		delChan: make(chan *Blip),
		getChan: make(chan chan []Blip),
		gw:      gw,
	}

	go func(bl BlipList) {
		for {
			select {
			case b := <-bl.addChan:
				bl.blips.PushBack(b)
			case bdel := <-bl.delChan:
				for e := bl.blips.Front(); e != nil; e = e.Next() {
					b, ok := e.Value.(*Blip)
					if !ok {
						fmt.Printf("Non-blip object in blip list!")
						os.Exit(2)
					}

					if bdel == b {
						bl.blips.Remove(e)
						break
					}
				}
			case bget := <-bl.getChan:
				bget <- bl.computeForOutput()
			}

		}
	}(bl)

	return bl
}

// add takes ownership over blip. It's forbidden from writing and reading.
func (bl *BlipList) Give(b *Blip) {
	bl.addChan <- b
}

func (bl *BlipList) Del(b *Blip) {
	bl.delChan <- b
}

func (bl *BlipList) Get() []Blip {
	rsp := make(chan []Blip)
	bl.getChan <- rsp
	return <-rsp
}

func (bl *BlipList) computeForOutput() []Blip {
	// List of blips that are being animated
	blipAnim := make([]*Blip, 0)
	// Keeps non-animated blips sorted by map location, for piling-up computation.
	blipMap := make([][]Blip, bl.gw.Tiles.Width*bl.gw.Tiles.Height)
	// Return list of blip copies.
	blipOutput := make([]Blip, 0, bl.blips.Len())

	for e := bl.blips.Front(); e != nil; e = e.Next() {
		b, ok := e.Value.(*Blip)
		if !ok {
			fmt.Printf("Non-blip object in blip list!")
			os.Exit(2)
		}

		if b.Path != nil {
			if b.Size == pixel.ZV {
				// Initialize animation
				b.Size = pixel.Vec{X: 2.0, Y: 2.0}
				b.animStart = time.Now()
			}
			blipAnim = append(blipAnim, b)
			continue
		}

		i := b.X + b.Y*bl.gw.Tiles.Width
		if blipMap[i] == nil {
			blipMap[i] = []Blip{*b}
		} else {
			blipMap[i] = append(blipMap[i], *b)
		}
	}

	// Update paint parameters of all stationary blips.
	for y := 0; y < bl.gw.Tiles.Height; y++ {
		for x := 0; x < bl.gw.Tiles.Width; x++ {
			tile := blipMap[x+y*bl.gw.Tiles.Width]
			if tile == nil {
				continue
			}

			basePos := bl.gw.TileToVec(x, y)
			// It was a bad idea to flip in gameWorld. Unflip.
			basePos.Y = float64(bl.gw.Tiles.Height*bl.gw.Tiles.TileHeight) - basePos.Y + 16.0

			basePos = basePos.Sub(pixel.Vec{
				X: float64(bl.gw.Tiles.TileWidth) / 2.0,
				Y: float64(bl.gw.Tiles.TileHeight) / 2.0,
			})

			fits := int(math.Sqrt(float64(len(tile))))
			if fits < 4 {
				fits = 4
			}
			// Assume square tile.
			size := int(float64(bl.gw.Tiles.TileWidth+1) / float64(fits))
			if size < 2 {
				size = 2
			}

			for i, b := range tile {
				dx := i % fits * size
				dy := i / fits * size
				b.Pos = basePos.Add(pixel.Vec{X: float64(dx), Y: float64(dy)})
				b.Size = pixel.Vec{X: float64(size) - 1.0, Y: float64(size) - 1.0}
				blipOutput = append(blipOutput, b)
			}
		}
	}

	// Update paint parameters of moving blips
	for _, ba := range blipAnim {
		currentT := float64(time.Since(ba.animStart) / time.Millisecond)
		var from, to pathNode
		// Remaining ms in this step (from->to)
		remainingT := 0.0

		pathT := 0.0
		for e := ba.Path.Front(); e != nil; e = e.Next() {
			pn, ok := e.Value.(pathNode)
			if !ok {
				fmt.Print("Non-pathNode found in path list\n")
				os.Exit(2)
			}

			if e == ba.Path.Front() {
				from = pn
				continue
			}

			pathT += pn.Cost()
			if pathT > currentT {
				// Found destination node
				to = pn
				remainingT = pathT - currentT
				break
			}

			from = pn
		}

		// At destination
		if to == nil {
			ba.Pos = bl.gw.TileToVec(from.X(), bl.gw.Tiles.Height-from.Y()-1)
		} else {
			progress := (to.Cost() - remainingT) / to.Cost()
			fromPos := bl.gw.TileToVec(from.X(), bl.gw.Tiles.Height-from.Y()-1)
			toPos := bl.gw.TileToVec(to.X(), bl.gw.Tiles.Height-to.Y()-1)
			d := toPos.Sub(fromPos).Scaled(progress)
			ba.Pos = fromPos.Add(d)
		}
		blipOutput = append(blipOutput, *ba)
	}

	return blipOutput
}
