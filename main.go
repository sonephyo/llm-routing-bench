package main

import (
	"fmt"
	"llm-routing-bench/backend"
)

func main() {
	ports := [...]string{"7000", "7001"}
	backends := []backend.Backend{}

	for _, port := range ports {
		backends = append(backends, backend.Backend{
			PortNumber: port,
			IsAlive: true,
		})
	}

	fmt.Println(backends)
}
