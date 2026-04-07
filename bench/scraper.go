package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type promResponse struct {
	Status string   `json:"status"`
	Data   promData `json:"data"`
}

type promData struct {
	Result []promResult `json:"result"`
}

type promResult struct {
	Metric map[string]string `json:"metric"`
	Value  []interface{}     `json:"value"` // [timestamp_float, value_string]
}

func queryProm(promAddr, query string) ([]promResult, error) {
	reqURL := promAddr + "/api/v1/query?query=" + url.QueryEscape(query)
	log.Printf("Address = %s", reqURL)
	resp, err := http.Get(reqURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var pr promResponse
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return nil, err
	}
	if pr.Status != "success" {
		return nil, fmt.Errorf("prometheus query failed: %s", pr.Status)
	}
	return pr.Data.Result, nil
}

func parseFloat(v []interface{}) float64 {
	if len(v) < 2 {
		return 0
	}
	s, ok := v[1].(string)
	if !ok {
		return 0
	}
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func querySimpleByBackend(promAddr, query string) (map[string]float64, error) {
	results, err := queryProm(promAddr, query)
	if err != nil {
		return nil, err
	}
	out := make(map[string]float64, len(results))
	for _, r := range results {
		backend := r.Metric["backend"]
		if backend == "" {
			backend = r.Metric["instance"]
		}
		out[backend] = parseFloat(r.Value)
	}
	return out, nil
}

func queryHistogramByBackend(promAddr, metricName string) (map[string]HistogramData, error) {
	buckets, err := queryProm(promAddr, metricName+"_bucket")
	if err != nil {
		return nil, err
	}
	sums, err := querySimpleByBackend(promAddr, metricName+"_sum")
	if err != nil {
		return nil, err
	}
	counts, err := querySimpleByBackend(promAddr, metricName+"_count")
	if err != nil {
		return nil, err
	}

	out := make(map[string]HistogramData)
	for _, r := range buckets {
		backend := r.Metric["backend"]
		if backend == "" {
			backend = r.Metric["instance"]
		}
		le := r.Metric["le"]
		h := out[backend]
		if h.Buckets == nil {
			h.Buckets = make(map[string]float64)
		}
		h.Buckets[le] = parseFloat(r.Value)
		out[backend] = h
	}
	for backend, s := range sums {
		h := out[backend]
		h.Sum = s
		out[backend] = h
	}
	for backend, c := range counts {
		h := out[backend]
		h.Count = c
		out[backend] = h
	}
	return out, nil
}


type routerSnapshot struct {
	RequestsTotal     map[string]float64
	DurationHistogram map[string]HistogramData
}

type backendSnapshot struct {
	NumRequestsWaiting    map[string]float64
	NumRequestsRunning    map[string]float64
	KVCacheUsagePerc      map[string]float64
	GenerationTokensTotal map[string]float64
	TTFTHistogram         map[string]HistogramData
	E2EHistogram          map[string]HistogramData
}

func ScrapeRouterMetrics(promAddr string) (routerSnapshot, error) {
	total, err := querySimpleByBackend(promAddr, "lb_requests_total")
	if err != nil {
		return routerSnapshot{}, fmt.Errorf("lb_requests_total: %w", err)
	}
	hist, err := queryHistogramByBackend(promAddr, "lb_request_duration_seconds")
	if err != nil {
		return routerSnapshot{}, fmt.Errorf("lb_request_duration_seconds: %w", err)
	}
	return routerSnapshot{RequestsTotal: total, DurationHistogram: hist}, nil
}

func ScrapeBackendMetrics(promAddr string) (*backendSnapshot, error) {
	waiting, err := querySimpleByBackend(promAddr, "vllm:num_requests_waiting")
	if err != nil || len(waiting) == 0 {
		return nil, nil
	}
	running, _ := querySimpleByBackend(promAddr, "vllm:num_requests_running")
	kv, _ := querySimpleByBackend(promAddr, "vllm:kv_cache_usage_perc")
	genTokens, _ := querySimpleByBackend(promAddr, "vllm:generation_tokens_total")
	ttft, _ := queryHistogramByBackend(promAddr, "vllm:time_to_first_token_seconds")
	e2e, _ := queryHistogramByBackend(promAddr, "vllm:e2e_request_latency_seconds")

	return &backendSnapshot{
		NumRequestsWaiting:    waiting,
		NumRequestsRunning:    running,
		KVCacheUsagePerc:      kv,
		GenerationTokensTotal: genTokens,
		TTFTHistogram:         ttft,
		E2EHistogram:          e2e,
	}, nil
}

func deltaSimpleMap(before, after map[string]float64) map[string]float64 {
	out := make(map[string]float64, len(after))
	for k, v := range after {
		out[k] = v - before[k]
	}
	return out
}

func deltaHistogramMap(before, after map[string]HistogramData) map[string]HistogramData {
	out := make(map[string]HistogramData, len(after))
	for backend, a := range after {
		b := before[backend]
		d := HistogramData{
			Buckets: make(map[string]float64, len(a.Buckets)),
			Sum:     a.Sum - b.Sum,
			Count:   a.Count - b.Count,
		}
		for le, cnt := range a.Buckets {
			d.Buckets[le] = cnt - b.Buckets[le]
		}
		out[backend] = d
	}
	return out
}

func gaugeSnapshotMap(before, after map[string]float64) map[string]GaugeSnapshot {
	out := make(map[string]GaugeSnapshot, len(after))
	for k, v := range after {
		out[k] = GaugeSnapshot{Start: before[k], End: v}
	}
	return out
}

func DeltaRouterMetrics(before, after routerSnapshot) RouterMetrics {
	return RouterMetrics{
		LBRequestsTotal:          deltaSimpleMap(before.RequestsTotal, after.RequestsTotal),
		LBRequestDurationSeconds: deltaHistogramMap(before.DurationHistogram, after.DurationHistogram),
	}
}

func waitForQueueDrain(promAddr string) {
	log.Printf("waiting for vLLM queues to drain...")
	for {
		running, err := querySimpleByBackend(promAddr, "vllm:num_requests_running")
		if err != nil || len(running) == 0 {
			return
		}
		waiting, _ := querySimpleByBackend(promAddr, "vllm:num_requests_waiting")

		total := 0.0
		for _, v := range running {
			total += v
		}
		for _, v := range waiting {
			total += v
		}

		if total == 0 {
			log.Printf("queues drained")
			return
		}
		log.Printf("queues not empty (%.0f requests in flight), waiting 5s...", total)
		time.Sleep(5 * time.Second)
	}
}

func DeltaBackendMetrics(before, after *backendSnapshot) *BackendMetrics {
	if before == nil || after == nil {
		return nil
	}
	return &BackendMetrics{
		NumRequestsWaiting:       gaugeSnapshotMap(before.NumRequestsWaiting, after.NumRequestsWaiting),
		NumRequestsRunning:       gaugeSnapshotMap(before.NumRequestsRunning, after.NumRequestsRunning),
		KVCacheUsagePerc:         gaugeSnapshotMap(before.KVCacheUsagePerc, after.KVCacheUsagePerc),
		GenerationTokensTotal:    deltaSimpleMap(before.GenerationTokensTotal, after.GenerationTokensTotal),
		TimeToFirstTokenSeconds:  deltaHistogramMap(before.TTFTHistogram, after.TTFTHistogram),
		E2ERequestLatencySeconds: deltaHistogramMap(before.E2EHistogram, after.E2EHistogram),
	}
}
