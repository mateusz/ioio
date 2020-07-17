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
	blips list.List
	add   chan blip
	del   chan blip
}

func NewBlipList() blipList {
	bl := blipList{
		add: make(chan blip),
		del: make(chan blip),
	}

	go func(bl blipList) {
		for {
			select {
			case b := <-bl.add:
				bl.blips.PushBack(b)
			case bdel := <-bl.del:
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
			}
		}
	}(bl)

	return bl
}
