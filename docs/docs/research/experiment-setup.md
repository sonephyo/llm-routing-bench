---
title: Experiment Setup & Methodology
sidebar_position: 3
---

A description of the test environment, workload design, and measurement approach used in this benchmark. This page will be updated as the experimental setup is finalized.

## Test Environment

_Hardware specs and model details to be filled in once finalized._

The benchmark runs multiple replicas of the same LLM inference server behind a routing layer. Each replica serves the same model. The routing layer is instrumented to switch between strategies without changing the inference servers themselves, ensuring a fair comparison.

## Workloads

_Workload details TBA._

The workloads are designed to stress-test routing strategies across different traffic patterns. Planned scenarios include:

- **Steady-state load** — a stable arrival rate sustained over time
- **Bursty traffic** — sudden spikes in request volume followed by quiet periods
- **Mixed request sizes** — a distribution of short and long prompts to expose cost variance between requests

The goal is to replicate conditions that approximate real production usage rather than synthetic uniform load, since uniform load tends to hide the differences between routing strategies.

## Ensuring Fairness

Each routing strategy is evaluated against the same sequence of requests and the same replica pool. The random seed for traffic generation is fixed across runs. Replicas are restarted to a clean state between strategy evaluations to prevent residual KV cache or queue state from carrying over.

_Additional controls TBA as methodology is finalized._

## Metrics Collected

The primary metric is **tail latency** — specifically p95 and p99 time-to-first-token (TTFT) and end-to-end latency. These are chosen because they reflect the experience of the slowest requests, which is where routing strategies diverge most visibly.

Secondary metrics include:

- **Throughput** — total requests completed per unit time
- **Queue depth per replica** — to understand how work distributes across the pool
- **KV cache utilization per replica** — relevant for the least-KV-cache strategy

Metrics are collected via Prometheus and visualized in Grafana. Raw data is retained for offline analysis.

## Reproducing the Experiment

_Full reproduction instructions TBA. At minimum the steps will cover:_

1. Provisioning the replica pool
2. Loading the model onto each replica
3. Configuring and starting the router with a given strategy
4. Running the traffic generator
5. Collecting and exporting metrics from Prometheus

The benchmark tooling and configuration files are available in the project repository.
