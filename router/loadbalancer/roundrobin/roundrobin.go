package roundrobin

import (
	"llm-routing-bench/router/backend"
	"sync"
)

type RoundRobin struct {
	curBackend int
	mu         sync.Mutex
}

func (rb *RoundRobin) Route(backends []backend.Backend) backend.Backend {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	idx := rb.curBackend % len(backends)
	rb.curBackend++
	return backends[idx]
}
