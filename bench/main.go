// bench/cmd/bench/main.go
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"net/http"
	"sort"
	"sync"
	"time"
)

type CompletionRequest struct {
	Model     string `json:"model"`
	Prompt    string `json:"prompt"`
	MaxTokens int    `json:"max_tokens"`
}

type Result struct {
	Latency   time.Duration
	Success   bool
	RequestID int
}

func main() {
	port       := flag.Int("port", 7999, "Target port")
	n          := flag.Int("n", 1000, "Number of requests")
	concurrency := flag.Int("c", 50, "Concurrent workers")
	strategy   := flag.String("strategy", "round-robin", "Routing strategy label")
	flag.Parse()

	experimentID := fmt.Sprintf("%s_%s", *strategy, time.Now().Format("20060102_150405"))
	url          := fmt.Sprintf("http://localhost:%d/v1/completions", *port)

	payload := CompletionRequest{
		Model:     "mistralai/Mistral-7B-v0.1",
		Prompt:    "The following is a detailed history of computer science from the 1940s through the present day, covering key innovations in hardware, software, networking, and artificial intelligence.",
		MaxTokens: 1000,
	}
	body, _ := json.Marshal(payload)

	fmt.Printf("Experiment : %s\n", experimentID)
	fmt.Printf("Target     : %s\n", url)
	fmt.Printf("Requests   : %d | Concurrency: %d\n\n", *n, *concurrency)

	results := make([]Result, *n)
	// sem     := make(chan struct{}, *concurrency)
	var wg  sync.WaitGroup

	startTime := time.Now()

	for i := 0; i < *n; i++ {
		wg.Add(1)
		// sem <- struct{}{} // block if C goroutines already running

		go func(idx int) {
			defer wg.Done()
			// defer func() { <-sem }()

			start := time.Now()
			resp, err := http.Post(url, "application/json", bytes.NewReader(body))
			elapsed := time.Since(start)

			success := err == nil && resp != nil && resp.StatusCode == 200
			if resp != nil {
				resp.Body.Close()
			}

			results[idx] = Result{
				Latency:   elapsed,
				Success:   success,
				RequestID: idx,
			}
		}(i)
	}

	wg.Wait()
	totalDuration := time.Since(startTime)

	// --- compute stats ---
	latencies := make([]float64, 0, *n)
	successCount := 0
	for _, r := range results {
		if r.Success {
			successCount++
			latencies = append(latencies, float64(r.Latency.Milliseconds()))
		}
	}
	sort.Float64s(latencies)

	p50 := percentile(latencies, 50)
	p95 := percentile(latencies, 95)
	p99 := percentile(latencies, 99)
	throughput := float64(*n) / totalDuration.Seconds()

	fmt.Println("=== Results ===")
	fmt.Printf("Experiment ID : %s\n", experimentID)
	fmt.Printf("Total time    : %v\n", totalDuration.Round(time.Millisecond))
	fmt.Printf("Throughput    : %.2f req/s\n", throughput)
	fmt.Printf("Success rate  : %d/%d\n", successCount, *n)
	fmt.Printf("p50 latency   : %.2f ms\n", p50)
	fmt.Printf("p95 latency   : %.2f ms\n", p95)
	fmt.Printf("p99 latency   : %.2f ms\n", p99)

	// --- save JSON ---
	summary := map[string]any{
		"experiment_id":    experimentID,
		"strategy":         *strategy,
		"num_requests":     *n,
		"concurrency":      *concurrency,
		"start_ts":         startTime.Unix(),
		"end_ts":           time.Now().Unix(),
		"duration_seconds": totalDuration.Seconds(),
		"throughput_rps":   throughput,
		"success_count":    successCount,
		"p50_ms":           p50,
		"p95_ms":           p95,
		"p99_ms":           p99,
	}
	out, _ := json.MarshalIndent(summary, "", "  ")
	fmt.Printf("\n%s\n", out)
}

func percentile(sorted []float64, p float64) float64 {
	if len(sorted) == 0 {
		return 0
	}
	idx := math.Ceil(float64(len(sorted))*p/100) - 1
	return sorted[int(idx)]
}