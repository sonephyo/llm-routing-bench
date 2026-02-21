package roundrobin

import (
	"fmt"
	"llm-routing-bench/router/backend"
)

type RoundRobin struct {
	curBackend int
}

func (rb *RoundRobin) Route(backends []backend.Backend) backend.Backend {
	fmt.Println("This is RoundRobin")
	idx := rb.curBackend % len(backends)
	rb.curBackend++
	return backends[idx]
}
