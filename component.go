package main

import "github.com/faiface/pixel"

type component struct {
	position pixel.Vec
	x        int
	y        int
	name     string
	con      string
	proc     int
	sched    string
	lat      string
	spriteID uint32
	cores    int
}
