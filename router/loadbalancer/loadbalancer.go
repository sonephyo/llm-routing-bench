package loadbalancer

import (
	"llm-routing-bench/router/backend"
)

type Router interface {
	Route() backend.Backend
}
