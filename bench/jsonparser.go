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
}

type HistogramData struct {
	Buckets map[string]float64 `json:"buckets"`
	Sum     float64            `json:"sum"`
	Count   float64            `json:"count"`
}

type GaugeSnapshot struct {
	Start float64 `json:"start"`
	End   float64 `json:"end"`
}

type RouterMetrics struct {
	LBRequestsTotal          map[string]float64       `json:"lb_requests_total"`
	LBRequestDurationSeconds map[string]HistogramData `json:"lb_request_duration_seconds"`
}

type BackendMetrics struct {
	NumRequestsWaiting       map[string]GaugeSnapshot `json:"vllm_num_requests_waiting"`
	NumRequestsRunning       map[string]GaugeSnapshot `json:"vllm_num_requests_running"`
	KVCacheUsagePerc         map[string]GaugeSnapshot `json:"vllm_kv_cache_usage_perc"`
	GenerationTokensTotal    map[string]float64       `json:"vllm_generation_tokens_total"`
	TimeToFirstTokenSeconds  map[string]HistogramData `json:"vllm_time_to_first_token_seconds"`
	E2ERequestLatencySeconds map[string]HistogramData `json:"vllm_e2e_request_latency_seconds"`
}

type ExperimentResult struct {
	Metadata       ExperimentMetadata `json:"metadata"`
	VegetaMetrics  vegeta.Metrics     `json:"vegeta_metrics"`
	RouterMetrics  RouterMetrics      `json:"router_metrics"`
	BackendMetrics *BackendMetrics    `json:"backend_metrics,omitempty"`
}
