package roundrobin

import (
	"llm-routing-bench/router/backend"
	"sync"
)

type RoundRobin struct {
	curBackend int
	mu         sync.Mutex
	backends   []backend.Backend
}

func NewRoundRobin(backends []backend.Backend) *RoundRobin {
	return &RoundRobin{
		backends: backends,
	}
}

func (rb *RoundRobin) Route() backend.Backend {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	idx := rb.curBackend % len(rb.backends)
	rb.curBackend++
	return rb.backends[idx]
}
