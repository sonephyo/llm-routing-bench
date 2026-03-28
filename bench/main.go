// bench/cmd/bench/main.go
package main

import (
	"flag"
	"fmt"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func main() {

	freq := flag.Int("freq", 1, "number of requests per second")
	seconds := flag.Int("duration", 10, "duration in seconds")

	rate := vegeta.Rate{Freq: *freq, Per: time.Second}
	duration := time.Duration(*seconds) * time.Second
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "GET",
		URL: "http://localhost:7999",
	})
	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Big Bang") {
		metrics.Add(res)
	}
	metrics.Close()

	fmt.Printf("99th percentile: %s\n", metrics.Latencies.P99)
	fmt.Printf("Throughput: %f\n", metrics.Throughput)
}
