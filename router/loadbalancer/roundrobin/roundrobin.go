package roundrobin

import (
	"fmt"
	"llm-routing-bench/router/backend"
)

type RoundRobin struct {
}

func (rb *RoundRobin) Route(backends []backend.Backend) backend.Backend {
	fmt.Println("This is RoundRobin")
	return backends[0]
}
