package leastkvcache

import (
	"llm-routing-bench/router/backend"
	"llm-routing-bench/router/scraper"
	"log"
	"math"
	"net/http"
	"sync"
	"time"
)

type LeastKVCache struct {
	backends []backend.Backend
	mu       sync.RWMutex
	cache    map[string]float64
}

func NewLeastKVCache(backends []backend.Backend, pollInterval time.Duration) *LeastKVCache {

	lkv := LeastKVCache{
		backends: backends,
		cache:    make(map[string]float64),
	}

	lkv.poll()

	go lkv.pollLoop(pollInterval)

	return &lkv
}

func (lkv *LeastKVCache) poll() {
	tempKVMap := make(map[string]float64)

	for i := range lkv.backends {
		metrics, err := scraper.GetFilteredMetrics(lkv.backends[i].BackendURI, []string{
			"vllm:kv_cache_usage_perc",
		})
		if err != nil {
			log.Printf("warn: failed to scrape %s: %v", lkv.backends[i].BackendURI, err)
			tempKVMap[lkv.backends[i].BackendURI] = math.Inf(1)
			continue
		}
		tempKVMap[lkv.backends[i].BackendURI] = metrics["vllm:kv_cache_usage_perc"]
	}

	lkv.mu.Lock()
	defer lkv.mu.Unlock()
	lkv.cache = tempKVMap
}

func (lkv *LeastKVCache) pollLoop(pollInterval time.Duration) {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()
	for range ticker.C {
		lkv.poll()
	}

}

func (lkv *LeastKVCache) Route(r *http.Request) *backend.Backend {
	if len(lkv.backends) == 0 {
		return nil
	}
	var selectedBackend *backend.Backend
	var minPercentage float64 = math.Inf(1)

	lkv.mu.RLock()
	defer lkv.mu.RUnlock()

	for i := range lkv.backends {
		kvcacheVal := lkv.cache[lkv.backends[i].BackendURI]
		if kvcacheVal < minPercentage {
			selectedBackend = &lkv.backends[i]
			minPercentage = kvcacheVal
		}
	}

	return selectedBackend
}
