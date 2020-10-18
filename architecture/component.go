package architecture

import "github.com/faiface/pixel"

// component represents an active element in the simulation, such as wire or cpu.
type Component struct {
	Position pixel.Vec
	X        int
	Y        int
	Name     string
	Con      string
	Proc     int
	Sched    string
	Lat      string
	SpriteID uint32
	Cores    int
	Queue    int
	Workers  int
}
