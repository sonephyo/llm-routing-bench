package leastqueue

import (
	"llm-routing-bench/router/backend"
	"llm-routing-bench/router/scraper"
	"log"
	"math"
	"net/http"
	"sync"
	"time"
)

type LeastQueue struct {
	backends []backend.Backend
	mu       sync.RWMutex
	cache    map[string]float64
}

func NewLeastQueue(backends []backend.Backend, pollInterval time.Duration) *LeastQueue {
	lq := &LeastQueue{
		backends: backends,
		cache:    make(map[string]float64),
	}

	lq.poll()

	go lq.pollLoop(pollInterval)

	return lq
}

func (lq *LeastQueue) pollLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		lq.poll()
	}
}

// Update lq.cache with lastest retrieved from vllm
func (lq *LeastQueue) poll() {
	tempLQMap := make(map[string]float64)
	for i := range lq.backends {
		metrics, err := scraper.GetFilteredMetrics(lq.backends[i].BackendURI, []string{
			"vllm:num_requests_running",
			"vllm:num_requests_waiting",
		})
		if err != nil {
			log.Printf("warn: failed to scrape %s: %v", lq.backends[i].BackendURI, err)
			tempLQMap[lq.backends[i].BackendURI] = math.Inf(1)
			continue
		}

		queueDepth := metrics["vllm:num_requests_running"] + metrics["vllm:num_requests_waiting"]
		tempLQMap[lq.backends[i].BackendURI] = queueDepth
	}

	lq.mu.Lock()
	defer lq.mu.Unlock()
	lq.cache = tempLQMap
}

func (lq *LeastQueue) Route(r *http.Request) *backend.Backend {
	if len(lq.backends) == 0 {
		return nil
	}
	var selectedServer *backend.Backend
	minVal := math.Inf(1)

	lq.mu.RLock()
	defer lq.mu.RUnlock()

	for i := range lq.backends {
		depth := lq.cache[lq.backends[i].BackendURI]
		if depth < minVal {
			selectedServer = &lq.backends[i]
			minVal = depth
		}
	}
	return selectedServer
}
