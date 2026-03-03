package loadbalancer

import (
	"llm-routing-bench/router/backend"
	"net/http"
)

type Router interface {
	Route(r *http.Request) *backend.Backend
}
