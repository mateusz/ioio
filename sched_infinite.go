package main

import (
	"log"
	"time"
)

type schedRequest struct {
	rsp chan bool
	c   int
}

type schedInfinite struct {
	roComponent  *component
	scheduleChan chan schedRequest
}

func NewSchedInfinite(c *component) *schedInfinite {
	s := &schedInfinite{
		roComponent:  c,
		scheduleChan: make(chan schedRequest),
	}
	go s.start()

	return s
}

func (s *schedInfinite) start() {
	for {
		r := <-s.scheduleChan
		go func(r schedRequest) {
			log.Printf("{%s} Scheduler consuming %d", s.roComponent.name, r.c)
			time.Sleep(time.Millisecond * time.Duration(r.c))
			r.rsp <- true
		}(r)
	}
}

// Infinite parallelism.
func (s *schedInfinite) schedule(c int) {
	req := schedRequest{
		c:   c,
		rsp: make(chan bool),
	}
	s.scheduleChan <- req
	<-req.rsp
}
