package leastkvcache

import (
	"llm-routing-bench/router/backend"
	"net/http"
	"sync"
	"time"
)

type LeastKVCache struct {
	backends []backend.Backend
	mu sync.RWMutex
	cache map[string]float64
}

func NewLeastKVCache(backends []backend.Backend, pollInterval time.Duration) *LeastKVCache {
	
	lkv := LeastKVCache{
		backends: backends,
		cache: make(map[string]float64),
	}

	lkv.poll()

	go lkv.pollInterval(pollInterval)
	

	return &lkv
}

func (lkv *LeastKVCache) poll() {
	// Operation related to polling
}

func (lkv *LeastKVCache) pollInterval(pollInterval time.Duration) {
	timer := time.NewTimer(pollInterval)
	go func() {
		<-timer.C
		lkv.poll()
	}()
}

func (lkv *LeastKVCache) Route(r *http.Request) *backend.Backend {
	if len(lkv.backends) == 0 {
		return nil
	}
	var selectedBackend *backend.Backend
	var minPercentage float64 = 101
	
	for i := range lkv.backends {
		kvcacheVal := lkv.cache[lkv.backends[i].BackendURI]
		if kvcacheVal < minPercentage {
			selectedBackend = &lkv.backends[i]
			minPercentage = kvcacheVal
		}
	}

	return selectedBackend
}