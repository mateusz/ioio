package main

import (
	"container/list"
	"fmt"
	"os"
	"time"
)

type schedProcess struct {
	req        schedRequest
	cRemaining int
}

type schedMultitasking struct {
	roComponent  *component
	scheduleChan chan schedRequest
	processess   *list.List
}

// Simulation of real-world core-limited machine
func NewSchedMultitasking(c *component) *schedMultitasking {
	s := &schedMultitasking{
		roComponent:  c,
		scheduleChan: make(chan schedRequest),
		processess:   list.New(),
	}
	go s.start()

	return s
}

func (s *schedMultitasking) start() {
	ticker := time.NewTicker(100 * time.Millisecond)
	for {
		select {
		case r := <-s.scheduleChan:
			s.processess.PushBack(&schedProcess{
				req:        r,
				cRemaining: r.c,
			})
		case <-ticker.C:
			for i := 0; i < s.roComponent.cores; i++ {
				s.consume(100)
			}
		}
	}
}

func (s *schedMultitasking) consume(c int) {
	curr := s.processess.Front()
	if curr == nil {
		return
	}

	p, ok := curr.Value.(*schedProcess)
	if !ok {
		fmt.Printf("[%s] Unexpected item in the processess queue\n", s.roComponent.name)
		os.Exit(2)
	}

	p.cRemaining -= 100
	if p.cRemaining < 0 {
		p.req.rsp <- true
		s.processess.Remove(curr)
	}

	s.processess.MoveToBack(curr)
}

func (s *schedMultitasking) schedule(c int) {
	req := schedRequest{
		c:   c,
		rsp: make(chan bool),
	}
	s.scheduleChan <- req
	<-req.rsp
}
