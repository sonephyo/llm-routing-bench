package leastqueue

import (
	"llm-routing-bench/router/backend"
	"llm-routing-bench/router/scraper"
	"log"
	"net/http"
)

type LeastQueue struct {
	backends []backend.Backend
}

func NewLeastQueue(backends []backend.Backend) *LeastQueue {
	return &LeastQueue{
		backends: backends,
	}
}

func (lq *LeastQueue) Route(r *http.Request) *backend.Backend {
	for i := range lq.backends {
		scraper.GetFilteredMetrics(lq.backends[i].BackendURI, []string{"vllm:num_requests_running"})
	}
	log.Println()
	return &lq.backends[0]
}
