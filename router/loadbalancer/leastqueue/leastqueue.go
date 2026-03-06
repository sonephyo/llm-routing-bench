package leastqueue

import (
	"llm-routing-bench/router/backend"
	"llm-routing-bench/router/scraper"
	"log"
	"math"
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
	if len(lq.backends) == 0 {
		return nil
	}
	var selectedServer *backend.Backend
	minVal := math.Inf(1)
	for i := range lq.backends {
		metrics, err := scraper.GetFilteredMetrics(lq.backends[i].BackendURI, []string{
			"vllm:num_requests_running",
			"vllm:num_requests_waiting",
		})
		if err != nil {
			log.Fatal(err)
		}

		queueDepth := metrics["vllm:num_requests_running"] + metrics["vllm:num_requests_waiting"]
		if queueDepth < minVal {
			log.Println(lq.backends[i].BackendURI + " - number of queue depth : ")
			log.Println(queueDepth)
			log.Println(".........")
			selectedServer = &lq.backends[i]
			minVal = queueDepth
		}

	}
	return selectedServer
}
