package random

import (
	"llm-routing-bench/router/backend"
	"math/rand"
	"net/http"
)

type Random struct {
	backends []backend.Backend
}

func NewRandom(backends []backend.Backend) *Random {
	return &Random{backends: backends}
}

func (rb *Random) Route(r *http.Request) *backend.Backend {
	if len(rb.backends) == 0 {
		return nil
	}
	idx := rand.Intn(len(rb.backends))
	return &rb.backends[idx]
}
