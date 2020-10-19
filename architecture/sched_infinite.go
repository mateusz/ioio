package architecture

import (
	"log"
	"time"
)

// schedInfinite is a simulation of a machine with infinite resources - e.g. The Internet
type SchedInfinite struct {
	roComponent  *Component
	scheduleChan chan schedRequest
}

func NewSchedInfinite(c *Component) *SchedInfinite {
	s := &SchedInfinite{
		roComponent:  c,
		scheduleChan: make(chan schedRequest),
	}
	go s.start()

	return s
}

func (s *SchedInfinite) start() {
	for {
		r := <-s.scheduleChan
		go func(r schedRequest) {
			log.Printf("{%s} Scheduler consuming %d", s.roComponent.Name, r.c)
			time.Sleep(r.c)
			r.rsp <- true
		}(r)
	}
}

func (s *SchedInfinite) Schedule(c time.Duration) {
	req := schedRequest{
		c:   c,
		rsp: make(chan bool),
	}
	s.scheduleChan <- req
	<-req.rsp
}
