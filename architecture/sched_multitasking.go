package architecture

import (
	"container/list"
	"fmt"
	"os"
	"time"
)

// SchedMultitasking is a simulation of core-limited machine
type SchedMultitasking struct {
	roComponent  *Component
	scheduleChan chan schedRequest
	processess   *list.List
}

// schedProcess is used to track execution progress of schedRequest
type schedProcess struct {
	req        schedRequest
	cRemaining time.Duration
}

func NewSchedMultitasking(c *Component) *SchedMultitasking {
	s := &SchedMultitasking{
		roComponent:  c,
		scheduleChan: make(chan schedRequest),
		processess:   list.New(),
	}
	go s.start()

	return s
}

func (s *SchedMultitasking) start() {
	ticker := time.NewTicker(100 * time.Millisecond)
	for {
		select {
		case r := <-s.scheduleChan:
			s.processess.PushBack(&schedProcess{
				req:        r,
				cRemaining: r.c,
			})
		case <-ticker.C:
			for i := 0; i < s.roComponent.Cores; i++ {
				s.consume(100)
			}
		}
	}
}

func (s *SchedMultitasking) consume(c int) {
	curr := s.processess.Front()
	if curr == nil {
		return
	}

	p, ok := curr.Value.(*schedProcess)
	if !ok {
		fmt.Printf("[%s] Unexpected item in the processess queue\n", s.roComponent.Name)
		os.Exit(2)
	}

	p.cRemaining -= 100
	if p.cRemaining < 0 {
		p.req.rsp <- true
		s.processess.Remove(curr)
	}

	s.processess.MoveToBack(curr)
}

func (s *SchedMultitasking) Schedule(c time.Duration) {
	req := schedRequest{
		c:   c,
		rsp: make(chan bool),
	}
	s.scheduleChan <- req
	<-req.rsp
}
