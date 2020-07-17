package main

type get struct {
	program  *program
	topLevel *topLevel
	ctl      ctl
	prg      prg
}

type ctl map[string]string
