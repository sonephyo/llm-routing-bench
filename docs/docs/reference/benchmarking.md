---
title: Benchmarking
sidebar_position: 2
---

The `bench/` tool sends concurrent requests to the router and reports latency percentiles. It is used to compare how routing strategies perform under load.

## Running the Benchmark

```bash
go run bench/main.go \
  -port 7999 \
  -n 1000 \
  -c 50 \
  -strategy roundrobin
```

## Flags

| Flag | Default | Description |
|---|---|---|
| `-port` | `7999` | Target router port |
| `-n` | `1000` | Total number of requests to send |
| `-c` | `50` | Number of concurrent workers |
| `-strategy` | `round-robin` | Label used in the experiment output (does not change routing behavior) |

## Output

The tool reports:

- **p50, p95, p99 latency** — the median and tail latency across all requests
- **Throughput** — requests per second
- **JSON summary** — machine-readable output of all metrics, written at the end of the run

The `-strategy` flag is a label only — it tags the JSON output so you can identify which run corresponds to which strategy when comparing results.
