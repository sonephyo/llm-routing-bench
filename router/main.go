package main

import (
	"bufio"
	"fmt"
	"llm-routing-bench/router/backend"
	"log"
	"net/http"
)

func main() {
	uri := "localhost"
	ports := [...]string{"8000", "8001"}
	backends := []backend.Backend{}

	for _, port := range ports {
		backends = append(backends, backend.Backend{
			PortNumber: port,
			IsAlive:    true,
		})
	}

	client := http.Client{}
	for _, backendPort := range backends {
		reqUrl := "http://" + uri + ":" + backendPort.PortNumber
		res, err := client.Get(reqUrl)
		if err != nil {
			log.Println("an error occurred:", err)
			panic(err)
		}

		fmt.Println("Response status:", res.Status)

		scanner := bufio.NewScanner(res.Body)
		for i := 0; scanner.Scan() && i < 5; i++ {

			fmt.Println("Response from " + backendPort.PortNumber + ":" + scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			panic(err)
		}
		res.Body.Close()
	}

}
