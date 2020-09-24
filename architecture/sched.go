package architecture

// schedRequest represents a computational request to any scheduler
type schedRequest struct {
	rsp chan bool
	c   int
}
