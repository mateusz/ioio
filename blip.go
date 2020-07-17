package main

import (
	"container/list"
	"fmt"
	"image/color"
	"os"
)

type blip struct {
	color color.Color
	x     int
	y     int
}

type blipList struct {
	blips   list.List
	addChan chan blip
	delChan chan blip
	getChan chan chan []blip
}

func NewBlipList() blipList {
	bl := blipList{
		addChan: make(chan blip),
		delChan: make(chan blip),
		getChan: make(chan chan []blip),
	}

	go func(bl blipList) {
		for {
			select {
			case b := <-bl.addChan:
				bl.blips.PushBack(b)
			case bdel := <-bl.delChan:
				for e := bl.blips.Front(); e != nil; e = e.Next() {
					b, ok := e.Value.(blip)
					if !ok {
						fmt.Printf("Non-blip object in blip list!")
						os.Exit(2)
					}

					if bdel.x == b.x && bdel.y == b.y && bdel.color == b.color {
						bl.blips.Remove(e)
						break
					}
				}
			case bget := <-bl.getChan:
				blips := make([]blip, 0, bl.blips.Len())
				for e := bl.blips.Front(); e != nil; e = e.Next() {
					b, ok := e.Value.(blip)
					if !ok {
						fmt.Printf("Non-blip object in blip list!")
						os.Exit(2)
					}

					blips = append(blips, b)
				}
				bget <- blips
			}

		}
	}(bl)

	return bl
}

func (bl *blipList) add(x, y int, color color.Color) {
	bl.addChan <- blip{
		x:     x,
		y:     y,
		color: color,
	}
}

func (bl *blipList) del(x, y int, color color.Color) {
	bl.delChan <- blip{
		x:     x,
		y:     y,
		color: color,
	}
}

func (bl *blipList) get() []blip {
	rsp := make(chan []blip)
	bl.getChan <- rsp
	return <-rsp
}
