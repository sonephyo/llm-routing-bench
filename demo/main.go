package main 

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func main() {

	freq := flag.Int("freq", 1, "number of requests per second")
	seconds := flag.Int("duration", 10, "duration in seconds")
	flag.Parse()

	rate := vegeta.Rate{Freq: *freq, Per: time.Second}
	duration := time.Duration(*seconds) * time.Second

	reqBody := `{"model": "mistralai/Mistral-7B-v0.1", "prompt": "The following is a detailed history of computer science from the 1940s through the present day, covering key innovations in hardware, software, networking, and artificial intelligence.", "max_tokens": 1000}`

	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "POST",
		URL:    "http://localhost:7999",
		Body:   []byte(reqBody),
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
	})
	attacker := vegeta.NewAttacker(vegeta.Timeout(300 * time.Second))

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Big Bang") {
		metrics.Add(res)
	}
	metrics.Close()

	fmt.Printf("99th percentile: %s\n", metrics.Latencies.P99)
	fmt.Printf("Throughput: %f\n", metrics.Throughput)
}
