package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"llm-routing-bench/router/backend"
	"llm-routing-bench/router/loadbalancer"
	"llm-routing-bench/router/loadbalancer/consistenthashing"
	"llm-routing-bench/router/loadbalancer/leastkvcache"
	"llm-routing-bench/router/loadbalancer/leastqueue"
	"llm-routing-bench/router/loadbalancer/random"
	"llm-routing-bench/router/loadbalancer/roundrobin"
	"llm-routing-bench/router/metrics"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	mode       string
	lbStrategy string
	rr         loadbalancer.Router
)

type LBServer struct {
	uri    string
	client http.Client
	router loadbalancer.Router
}

type ServerResponse struct {
	Message  string `json:"message"`
	Status   string `json:"status"`
	Response string `json:"response"`
}

func (lb *LBServer) backendHandler(w http.ResponseWriter, r *http.Request) {

	// Select backend to send out request
	selectedBackend := lb.router.Route(r)
	if selectedBackend == nil {
		http.Error(w, "no backend available", http.StatusServiceUnavailable)
		return
	}

	metrics.RequestCount.WithLabelValues(selectedBackend.BackendURI).Inc()
	start := time.Now()
	defer func() {
		metrics.RequestLatency.WithLabelValues(selectedBackend.BackendURI).Observe(time.Since(start).Seconds())
	}()

	w.Header().Set("Content-Type", "application/json")

	switch mode {
	// Local mode uses fake servers (HTTP servers) on local machine
	case "local":
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			response := ServerResponse{
				Message:  "Request not supported",
				Status:   "Error",
				Response: "",
			}

			err := json.NewEncoder(w).Encode(response)
			if err != nil {
				log.Printf("error encoding response: %v", err)
				return
			}
			return
		}
		resp, err := http.Get(selectedBackend.BackendURI)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		response := ServerResponse{
			Message:  "Selected Port Number: " + selectedBackend.BackendURI,
			Status:   "OK",
			Response: string(body),
		}

		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Printf("error encoding response: %v", err)
			return
		}
	// Server mode uses vllm servers as defined in docker-compose.yml
	case "server":
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			response := ServerResponse{
				Message:  "Request not supported",
				Status:   "Error",
				Response: "",
			}

			err := json.NewEncoder(w).Encode(response)
			if err != nil {
				log.Printf("error encoding response: %v", err)
				return
			}
			return
		}
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resp, err := http.Post(selectedBackend.BackendURI+"/v1/completions", "application/json", bytes.NewReader(bodyBytes))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		response := ServerResponse{
			Message:  "Selected Port Number: " + selectedBackend.BackendURI,
			Status:   "OK",
			Response: string(body),
		}

		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Printf("error encoding response: %v", err)
			return
		}

	}
}

func main() {

	mode = os.Getenv("MODE")
	fmt.Println("Selected mode: ", mode)
	lbStrategy = os.Getenv("LB_STRATEGY")
	fmt.Println("Selected strategy: ", lbStrategy)

	ports := [...]string{"http://backend-1:8000", "http://backend-2:8000"}
	backends := []backend.Backend{}

	for _, port := range ports {
		backends = append(backends, backend.Backend{
			BackendURI: port,
			IsAlive:    true,
		})
	}

	switch lbStrategy {
	case "random":
		rr = random.NewRandom(backends)
	case "roundrobin":
		rr = roundrobin.NewRoundRobin(backends)
	case "consistenthashing":
		rr = consistenthashing.NewConsistentHash(backends)
	case "leastqueue":
		rr = leastqueue.NewLeastQueue(backends, 250*time.Millisecond)
	case "leastkvcache":
		rr = leastkvcache.NewLeastKVCache(backends, 250*time.Millisecond)
	default:
		log.Fatalln("Invalid lbStrategy")
	}

	lbserver := LBServer{
		router: rr,
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

	log.Fatal(http.ListenAndServe(":7999", nil))
}
