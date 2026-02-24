package main

import (
	"encoding/json"
	"fmt"
	"llm-routing-bench/router/backend"
	"llm-routing-bench/router/loadbalancer"
	"llm-routing-bench/router/loadbalancer/roundrobin"
	"llm-routing-bench/router/metrics"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type LBServer struct {
	uri      string
	backends []backend.Backend
	client   http.Client
	router   loadbalancer.Router
}

type TempResponse struct {
	Message string
	Status  string
}

func (lb *LBServer) backendHandler(w http.ResponseWriter, r *http.Request) {

	selectedBackend := lb.router.Route(lb.backends)
	metrics.RequestCount.WithLabelValues(selectedBackend.PortNumber).Inc()
	fmt.Println(selectedBackend)
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	response := TempResponse{
		Message: "Selected Port Number: " + selectedBackend.PortNumber,
		Status:  "OK",
	}

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("error encoding response: %v", err)
		return
	}
}

func main() {
	uri := "localhost"
	ports := [...]string{"8000", "8001"}
	backends := []backend.Backend{}
	client := http.Client{}

	for _, port := range ports {
		backends = append(backends, backend.Backend{
			PortNumber: port,
			IsAlive:    true,
		})
	}

	rr := &roundrobin.RoundRobin{}

	lbserver := LBServer{
		uri:      uri,
		backends: backends,
		client:   client,
		router:   rr,
	}

	// Prometheus Related Endp
	reg := prometheus.NewRegistry()
	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		metrics.RequestCount,
		metrics.RequestLatency,
	)
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	http.HandleFunc("/", lbserver.backendHandler)

	http.ListenAndServe(":7999", nil)
}
