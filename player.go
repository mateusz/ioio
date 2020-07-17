package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type player struct {
	position      pixel.Vec
	scrollSpeed   float64
	scrollHotZone float64
	scrolling     pixel.Vec
}

func (p *player) Update(dt float64) {
	p.position = p.position.Add(p.scrolling.Scaled(dt))
}

func (p *player) Input(win *pixelgl.Window, cam pixel.Matrix) {
	p.scrolling = pixel.ZV
	if win.MouseInsideWindow() {
		if win.MousePosition().X < p.scrollHotZone {
			p.scrolling.X = -p.scrollSpeed
		}
		if win.MousePosition().X > monW-p.scrollHotZone {
			p.scrolling.X = p.scrollSpeed
		}
		if win.MousePosition().Y < p.scrollHotZone {
			p.scrolling.Y = -p.scrollSpeed
		}
		if win.MousePosition().Y > monH-p.scrollHotZone {
			p.scrolling.Y = p.scrollSpeed
		}
	}
}
