package main

import (
	"fmt"
	"llm-routing-bench/router/backend"
	"llm-routing-bench/router/loadbalancer"
	"llm-routing-bench/router/loadbalancer/roundrobin"
	"net/http"
)

type LBServer struct {
	uri      string
	backends []backend.Backend
	client   http.Client
	router   loadbalancer.Router
}

func (lb *LBServer) backendHandler(w http.ResponseWriter, r *http.Request) {

	selectedBackend := lb.router.Route(lb.backends)
	fmt.Println(selectedBackend)
	// for _, backendPort := range lb.backends {
	// 	reqUrl := "http://" + lb.uri + ":" + backendPort.PortNumber
	// 	res, err := lb.client.Get(reqUrl)
	// 	if err != nil {
	// 		http.Error(w, "backend error", http.StatusBadGateway)
	// 		return
	// 	}

	// 	fmt.Println("Response status:", res.Status)

	// 	body, err := io.ReadAll(res.Body)
	// 	fmt.Println(string(body))
	// 	res.Body.Close()
	// }
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

	http.HandleFunc("/", lbserver.backendHandler)

	http.ListenAndServe(":7999", nil)

}
