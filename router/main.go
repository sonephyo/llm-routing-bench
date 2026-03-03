package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"llm-routing-bench/router/backend"
	"llm-routing-bench/router/loadbalancer"
	"llm-routing-bench/router/loadbalancer/roundrobin"
	"llm-routing-bench/router/metrics"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	mode string
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

type vllmBackendStruct struct {
	Model     string `json:"model"`
	Prompt    string `json:"prompt"`
	MaxTokens int64  `json:"max_tokens"`
}

func (lb *LBServer) backendHandler(w http.ResponseWriter, r *http.Request) {

	selectedBackend := lb.router.Route()
	metrics.RequestCount.WithLabelValues(selectedBackend.BackendURI).Inc()
	fmt.Println(selectedBackend)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	switch mode {
	case "local":
		if r.Method != "GET" {
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
			log.Fatalln(err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
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
	case "server":
		if r.Method != "POST" {
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
			log.Fatalln(err)
		}

		fmt.Println("Body bytes:", string(bodyBytes))

		resp, err := http.Post(selectedBackend.BackendURI+"/v1/completions", "application/json", bytes.NewReader(bodyBytes))
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
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

	uri := "localhost"
	ports := [...]string{"http://backend-1:8000", "http://backend-2:8000"}
	backends := []backend.Backend{}
	client := http.Client{}

	for _, port := range ports {
		backends = append(backends, backend.Backend{
			BackendURI: port,
			IsAlive:    true,
		})
	}

	rr := roundrobin.NewRoundRobin(backends)

	lbserver := LBServer{
		uri:    uri,
		client: client,
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

	http.ListenAndServe(":7999", nil)
}
