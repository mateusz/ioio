package graphics

import (
	"container/list"
	"fmt"
	"image/color"
	"log"
	"math"
	"os"
	"time"

	"github.com/faiface/pixel"
	"github.com/mateusz/ioio/pathfinder"
	"github.com/mateusz/rtsian/piksele"
)

type Blip struct {
	Color     color.Color
	X         int        // logical location
	Y         int        // logcial location
	Path      *list.List // current movement path to target
	Size      pixel.Vec  // pain size
	Pos       pixel.Vec  // paint location
	animStart time.Time  // start of current movement target
	target    pixel.Vec  // position of current movement target
	d         float64    // distance to current movement target
	v         pixel.Vec  // velocity vector
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
				// It was a bad idea to flip in gameWorld. Unflip.
				b.Pos = bl.gw.TileToVec(b.X, bl.gw.Tiles.Height-b.Y-1)
				b.Size = pixel.Vec{X: 2.0, Y: 2.0}
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

		if ba.d > 0.0 {
			dt := time.Since(ba.animStart)
			v := ba.v.Scaled(float64(dt) / float64(time.Second))
			ba.Pos = ba.Pos.Add(v)
			ba.d -= v.Len()
			ba.animStart = time.Now()
			if ba.d < 0.0 {
				ba.Pos = ba.target
			}
		} else {
			bl.applyPath(ba)
		}

		blipOutput = append(blipOutput, *ba)
	}

	return blipOutput
}

func (bl *BlipList) applyPath(b *Blip) {
	if b.Path == nil || b.Path.Len() == 0 {
		b.target = b.Pos
		b.d = 0.0
		b.v = pixel.ZV
		b.Path = nil
		return
	}

	// Next path step
	n, ok := b.Path.Remove(b.Path.Front()).(*pathfinder.PathNode)
	if !ok {
		log.Panic("Fatal: path list contained non-pathNode!")
	}

	b.animStart = time.Now()
	// It was a bad idea to flip in.gameWorld. Unflip.
	b.target = bl.gw.TileToVec(n.X, bl.gw.Tiles.Height-n.Y-1)
	mv := b.target.Sub(b.Pos)
	b.d = mv.Len()
	b.v = mv.Unit().Scaled(float64(bl.gw.Tiles.Width) * 1000.0 / n.Cost)
}
