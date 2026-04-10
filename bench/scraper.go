package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
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

// promRangeResponse is the Prometheus /api/v1/query_range response envelope.
type promRangeResponse struct {
	Status string        `json:"status"`
	Data   promRangeData `json:"data"`
}

type promRangeData struct {
	Result []promRangeResult `json:"result"`
}

type promRangeResult struct {
	Metric map[string]string `json:"metric"`
	Values [][]interface{}   `json:"values"` // [[timestamp_float, value_string], ...]
}

var promClient = &http.Client{Timeout: 10 * time.Second}

func queryPromRange(promAddr, query string, start, end time.Time, step time.Duration) ([]promRangeResult, error) {
	params := url.Values{}
	params.Set("query", query)
	params.Set("start", fmt.Sprintf("%d", start.Unix()))
	params.Set("end", fmt.Sprintf("%d", end.Unix()))
	params.Set("step", fmt.Sprintf("%ds", int(step.Seconds())))
	reqURL := promAddr + "/api/v1/query_range?" + params.Encode()

	resp, err := promClient.Get(reqURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var pr promRangeResponse
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return nil, err
	}
	if pr.Status != "success" {
		return nil, fmt.Errorf("prometheus range query failed: %s", pr.Status)
	}
	return pr.Data.Result, nil
}

// scrapeGaugeRange fetches a gauge metric over [start, end] and returns
// min/max/mean per backend derived from the 1s-resolution time series.
func scrapeGaugeRange(promAddr, metric string, start, end time.Time) map[string]GaugeSeries {
	results, err := queryPromRange(promAddr, metric, start, end, time.Second)
	if err != nil {
		log.Printf("warn: range query %s: %v", metric, err)
		return nil
	}
	out := make(map[string]GaugeSeries, len(results))
	for _, r := range results {
		backend := r.Metric["backend"]
		if backend == "" {
			backend = r.Metric["instance"]
		}
		var sum float64
		mn := math.MaxFloat64
		mx := -math.MaxFloat64
		n := 0
		for _, v := range r.Values {
			if len(v) < 2 {
				continue
			}
			s, ok := v[1].(string)
			if !ok {
				continue
			}
			f, err := strconv.ParseFloat(s, 64)
			if err != nil {
				continue
			}
			sum += f
			if f < mn {
				mn = f
			}
			if f > mx {
				mx = f
			}
			n++
		}
		if n == 0 {
			out[backend] = GaugeSeries{}
			continue
		}
		out[backend] = GaugeSeries{Min: mn, Max: mx, Mean: sum / float64(n)}
	}
	return out
}

// gaugeRangeSnapshot holds min/max/mean time-series data for gauge metrics
// captured via range queries over the experiment window.
type gaugeRangeSnapshot struct {
	NumRequestsWaiting map[string]GaugeSeries
	NumRequestsRunning map[string]GaugeSeries
	KVCacheUsagePerc   map[string]GaugeSeries
}

func ScrapeGaugeRanges(promAddr string, start, end time.Time) *gaugeRangeSnapshot {
	return &gaugeRangeSnapshot{
		NumRequestsWaiting: scrapeGaugeRange(promAddr, "vllm:num_requests_waiting", start, end),
		NumRequestsRunning: scrapeGaugeRange(promAddr, "vllm:num_requests_running", start, end),
		KVCacheUsagePerc:   scrapeGaugeRange(promAddr, "vllm:kv_cache_usage_perc", start, end),
	}
}

func queryProm(promAddr, query string) ([]promResult, error) {
	reqURL := promAddr + "/api/v1/query?query=" + url.QueryEscape(query)
	log.Printf("Address = %s", reqURL)
	resp, err := promClient.Get(reqURL)
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
	GenerationTokensTotal map[string]float64
	PreemptionsTotal      map[string]float64
	PrefixCacheHits       map[string]float64
	PrefixCacheQueries    map[string]float64
	TTFTHistogram         map[string]HistogramData
	E2EHistogram          map[string]HistogramData
	QueueTimeHistogram    map[string]HistogramData
	PrefillTimeHistogram  map[string]HistogramData
	DecodeTimeHistogram   map[string]HistogramData
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
	// Use num_requests_waiting as a probe — if absent, backends aren't up.
	probe, err := querySimpleByBackend(promAddr, "vllm:num_requests_waiting")
	if err != nil || len(probe) == 0 {
		return nil, nil
	}
	genTokens, _ := querySimpleByBackend(promAddr, "vllm:generation_tokens_total")
	preemptions, _ := querySimpleByBackend(promAddr, "vllm:num_preemptions_total")
	prefixHits, _ := querySimpleByBackend(promAddr, "vllm:prefix_cache_hits_total")
	prefixQueries, _ := querySimpleByBackend(promAddr, "vllm:prefix_cache_queries_total")
	ttft, _ := queryHistogramByBackend(promAddr, "vllm:time_to_first_token_seconds")
	e2e, _ := queryHistogramByBackend(promAddr, "vllm:e2e_request_latency_seconds")
	queueTime, _ := queryHistogramByBackend(promAddr, "vllm:request_queue_time_seconds")
	prefillTime, _ := queryHistogramByBackend(promAddr, "vllm:request_prefill_time_seconds")
	decodeTime, _ := queryHistogramByBackend(promAddr, "vllm:request_decode_time_seconds")

	return &backendSnapshot{
		GenerationTokensTotal: genTokens,
		PreemptionsTotal:      preemptions,
		PrefixCacheHits:       prefixHits,
		PrefixCacheQueries:    prefixQueries,
		TTFTHistogram:         ttft,
		E2EHistogram:          e2e,
		QueueTimeHistogram:    queueTime,
		PrefillTimeHistogram:  prefillTime,
		DecodeTimeHistogram:   decodeTime,
	}, nil
}

func deltaSimpleMap(before, after map[string]float64) map[string]float64 {
	out := make(map[string]float64, len(after))
	for k, v := range after {
		if delta := v - before[k]; delta >= 0 {
			out[k] = delta
		} else {
			// Counter reset (e.g. container restarted) — after value is already
			// relative to the reset, so use it directly.
			out[k] = v
		}
	}
	return out
}

func deltaHistogramMap(before, after map[string]HistogramData) map[string]HistogramData {
	out := make(map[string]HistogramData, len(after))
	for backend, a := range after {
		b := before[backend]
		// If count went backwards the counter reset — use after values directly.
		reset := a.Count-b.Count < 0
		d := HistogramData{
			Buckets: make(map[string]float64, len(a.Buckets)),
		}
		if reset {
			d.Sum = a.Sum
			d.Count = a.Count
		} else {
			d.Sum = a.Sum - b.Sum
			d.Count = a.Count - b.Count
		}
		for le, cnt := range a.Buckets {
			if reset {
				d.Buckets[le] = cnt
			} else {
				d.Buckets[le] = cnt - b.Buckets[le]
			}
		}
		out[backend] = d
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

func DeltaBackendMetrics(before, after *backendSnapshot, gauges *gaugeRangeSnapshot) *BackendMetrics {
	if before == nil || after == nil {
		return nil
	}
	bm := &BackendMetrics{
		GenerationTokensTotal:     deltaSimpleMap(before.GenerationTokensTotal, after.GenerationTokensTotal),
		NumPreemptions:            deltaSimpleMap(before.PreemptionsTotal, after.PreemptionsTotal),
		PrefixCacheHits:           deltaSimpleMap(before.PrefixCacheHits, after.PrefixCacheHits),
		PrefixCacheQueries:        deltaSimpleMap(before.PrefixCacheQueries, after.PrefixCacheQueries),
		TimeToFirstTokenSeconds:   deltaHistogramMap(before.TTFTHistogram, after.TTFTHistogram),
		E2ERequestLatencySeconds:  deltaHistogramMap(before.E2EHistogram, after.E2EHistogram),
		RequestQueueTimeSeconds:   deltaHistogramMap(before.QueueTimeHistogram, after.QueueTimeHistogram),
		RequestPrefillTimeSeconds: deltaHistogramMap(before.PrefillTimeHistogram, after.PrefillTimeHistogram),
		RequestDecodeTimeSeconds:  deltaHistogramMap(before.DecodeTimeHistogram, after.DecodeTimeHistogram),
	}
	if gauges != nil {
		bm.NumRequestsWaiting = gauges.NumRequestsWaiting
		bm.NumRequestsRunning = gauges.NumRequestsRunning
		bm.KVCacheUsagePerc = gauges.KVCacheUsagePerc
	}
	return bm
}
