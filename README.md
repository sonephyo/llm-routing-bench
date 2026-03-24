# llm-routing-bench

A Go-based load balancer for LLM inference servers (vLLM) with pluggable routing strategies and integrated benchmarking. The goal is to measure how different routing strategies affect tail latency (p95/p99) in LLM inference serving.

## Architecture

```
Client/Benchmark Program
  |
  v
Router (Go, :7999)          <-- load balancer + Prometheus metrics
  |           |
  v           v
backend-1   backend-2       <-- vLLM servers (or fake servers in dev mode)
(Mistral-7B) (Mistral-7B)

Prometheus (:7779 / :7998)  <-- scrapes router /metrics
Grafana    (:7780 / :8000)  <-- dashboards
```

## Routing Strategies

| Strategy | Flag | Description |
|---|---|---|
| Round Robin | `roundrobin` | Cycles through backends sequentially |
| Consistent Hashing | `consistanthashing` | FNV-32a hash of request URL maps to a backend on a consistent ring |
| Least Queue | `leastqueue` | Scrapes `vllm:num_requests_running` from each backend and routes to the least loaded (WIP) |

## Modes

The router supports two modes, set via the `MODE` env var:

- **`server`** — Proxies POST requests to backend `/v1/completions` (production, real vLLM)
- **`local`** — Proxies GET requests to backends (dev/testing with fake servers)

## Quick Start

### Prerequisites

- Docker + Docker Compose
- Copy `.env.example` to `.env` and fill in values

```bash
cp .env.example .env
```

`.env` fields:

```
MODE=server           # or local
LB_STRATEGY=roundrobin  # roundrobin | consistanthashing | leastqueue
```

### Dev Mode (no GPU required)

Uses fake backend servers that return a simple JSON response.

```bash
# Start stack (pass router strategy as make variable)
make local-up router=roundrobin

# Stop stack
make local-down
```

Dev ports:

| Service | Port |
|---|---|
| Backend 1 | 7777 |
| Backend 2 | 7778 |
| Router | 7999 |
| Prometheus | 7998 |
| Grafana | 8000 |

### Production Mode (NVIDIA GPU required)

Runs two vLLM instances serving `mistralai/Mistral-7B-v0.1`, each pinned to a separate GPU.

```bash
make server-up router=roundrobin

make server-down
```

Production ports:

| Service | Port |
|---|---|
| Router | 7999 |
| Prometheus | 7779 |
| Grafana | 7780 |

### Health Check

```bash
make health
```

## Benchmarking

The `bench/` tool sends concurrent requests to the router and reports latency percentiles.

```bash
go run bench/main.go \
  -port 7999 \
  -n 1000 \
  -c 50 \
  -strategy roundrobin
```

Flags:

| Flag | Default | Description |
|---|---|---|
| `-port` | 7999 | Target router port |
| `-n` | 1000 | Total number of requests |
| `-c` | 50 | Concurrent workers |
| `-strategy` | round-robin | Label for the experiment ID |

Output includes p50, p95, p99 latency, throughput (req/s), and a JSON summary.

## Metrics

The router exposes Prometheus metrics at `:7999/metrics`.

| Metric | Type | Description |
|---|---|---|
| `lb_requests_total` | Counter | Total requests routed, labeled by backend |
| `lb_request_duration_seconds` | Histogram | Request latency, labeled by backend |

## Makefile Reference

| Target | Description |
|---|---|
| `make server-up router=<strategy>` | Build and start production stack |
| `make server-down` | Stop production stack |
| `make local-up router=<strategy>` | Build and start dev stack |
| `make local-down` | Stop dev stack |
| `make health` | Check health of all services |
| `make local-test` | Send test requests directly to backends |
| `make logs` | Follow Docker Compose logs |
