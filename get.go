package main

import "image/color"

type get struct {
	program  *program
	topLevel *topLevel
	ctl      ctl
	prg      prg
	color    *color.Color
}

type ctl map[string]string
