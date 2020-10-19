package architecture

import "time"

// schedRequest represents a computational request to any scheduler
type schedRequest struct {
	rsp chan bool
	c   time.Duration
}
