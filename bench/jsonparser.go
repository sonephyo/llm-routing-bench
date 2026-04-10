package main

import (
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type ExperimentMetadata struct {
	Strategy    string    `json:"strategy"`
	LoadPattern string    `json:"load_pattern"`
	TokenSize   int       `json:"token_size"`
	PromptType  string    `json:"prompt_type"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
}

type HistogramData struct {
	Buckets map[string]float64 `json:"buckets"`
	Sum     float64            `json:"sum"`
	Count   float64            `json:"count"`
}

// GaugeSeries captures the min/max/mean of a gauge metric over the experiment
// window via Prometheus range queries, replacing the meaningless start/end snapshot.
type GaugeSeries struct {
	Min  float64 `json:"min"`
	Max  float64 `json:"max"`
	Mean float64 `json:"mean"`
}

type RouterMetrics struct {
	LBRequestsTotal          map[string]float64       `json:"lb_requests_total"`
	LBRequestDurationSeconds map[string]HistogramData `json:"lb_request_duration_seconds"`
}

type BackendMetrics struct {
	NumRequestsWaiting        map[string]GaugeSeries   `json:"vllm_num_requests_waiting"`
	NumRequestsRunning        map[string]GaugeSeries   `json:"vllm_num_requests_running"`
	KVCacheUsagePerc          map[string]GaugeSeries   `json:"vllm_kv_cache_usage_perc"`
	GenerationTokensTotal     map[string]float64       `json:"vllm_generation_tokens_total"`
	NumPreemptions            map[string]float64       `json:"vllm_num_preemptions"`
	PrefixCacheHits           map[string]float64       `json:"vllm_prefix_cache_hits"`
	PrefixCacheQueries        map[string]float64       `json:"vllm_prefix_cache_queries"`
	TimeToFirstTokenSeconds   map[string]HistogramData `json:"vllm_time_to_first_token_seconds"`
	E2ERequestLatencySeconds  map[string]HistogramData `json:"vllm_e2e_request_latency_seconds"`
	RequestQueueTimeSeconds   map[string]HistogramData `json:"vllm_request_queue_time_seconds"`
	RequestPrefillTimeSeconds map[string]HistogramData `json:"vllm_request_prefill_time_seconds"`
	RequestDecodeTimeSeconds  map[string]HistogramData `json:"vllm_request_decode_time_seconds"`
}

// RawRequest captures per-request data for phase-level latency analysis.
// TimestampNs is Unix nanoseconds when the request was dispatched;
// LatencyNs is the end-to-end response time in nanoseconds.
type RawRequest struct {
	TimestampNs int64 `json:"timestamp_ns"`
	LatencyNs   int64 `json:"latency_ns"`
}

type ExperimentResult struct {
	Metadata       ExperimentMetadata `json:"metadata"`
	VegetaMetrics  vegeta.Metrics     `json:"vegeta_metrics"`
	RouterMetrics  RouterMetrics      `json:"router_metrics"`
	BackendMetrics *BackendMetrics    `json:"backend_metrics,omitempty"`
	RawRequests    []RawRequest       `json:"raw_requests"`
}
