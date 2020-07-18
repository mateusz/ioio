package main

import (
	"container/list"
	"fmt"
	"image/color"
	"math"
	"os"

	"github.com/faiface/pixel"
)

type blip struct {
	color color.Color
	x     int
	y     int
	pos   pixel.Vec
	size  pixel.Vec
}

type blipList struct {
	blips   list.List
	addChan chan *blip
	delChan chan *blip
	getChan chan chan []blip
}

func NewBlipList() blipList {
	bl := blipList{
		addChan: make(chan *blip),
		delChan: make(chan *blip),
		getChan: make(chan chan []blip),
	}

	go func(bl blipList) {
		for {
			select {
			case b := <-bl.addChan:
				bl.blips.PushBack(b)
			case bdel := <-bl.delChan:
				for e := bl.blips.Front(); e != nil; e = e.Next() {
					b, ok := e.Value.(*blip)
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
func (bl *blipList) add(b *blip) {
	bl.addChan <- b
}

func (bl *blipList) del(b *blip) {
	bl.delChan <- b
}

func (bl *blipList) get() []blip {
	rsp := make(chan []blip)
	bl.getChan <- rsp
	return <-rsp
}

func (bl *blipList) computeForOutput() []blip {
	blipMap := make([][]blip, gameWorld.Tiles.Width*gameWorld.Tiles.Height)
	blipOutput := make([]blip, 0, bl.blips.Len())

	for e := bl.blips.Front(); e != nil; e = e.Next() {
		b, ok := e.Value.(*blip)
		if !ok {
			fmt.Printf("Non-blip object in blip list!")
			os.Exit(2)
		}

		i := b.x + b.y*gameWorld.Tiles.Width
		if blipMap[i] == nil {
			blipMap[i] = []blip{*b}
		} else {
			blipMap[i] = append(blipMap[i], *b)
		}
	}

	for y := 0; y < gameWorld.Tiles.Height; y++ {
		for x := 0; x < gameWorld.Tiles.Width; x++ {
			tile := blipMap[x+y*gameWorld.Tiles.Width]
			if tile == nil {
				continue
			}
			basePos := gameWorld.TileToVec(x, y)
			basePos = basePos.Sub(pixel.Vec{
				X: float64(gameWorld.Tiles.TileWidth) / 2.0,
				Y: float64(gameWorld.Tiles.TileHeight) / 2.0,
			})

			fits := int(math.Sqrt(float64(len(tile))))
			if fits < 4 {
				fits = 4
			}
			// Assume square tile.
			size := int(float64(gameWorld.Tiles.TileWidth+1) / float64(fits))

			for i, b := range tile {
				dx := i % fits * size
				dy := i / fits * size
				b.pos = basePos.Add(pixel.Vec{X: float64(dx), Y: float64(dy)})
				b.size = pixel.Vec{X: float64(size) - 1.0, Y: float64(size) - 1.0}
				blipOutput = append(blipOutput, b)
			}
		}
	}

	return blipOutput
}
