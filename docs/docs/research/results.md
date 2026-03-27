---
title: Results & Analysis
sidebar_position: 4
---

Findings from the benchmark comparing round-robin, consistent hashing, least-queue, and least-KV-cache routing strategies under production-like LLM inference workloads. This page will be updated as experiments complete.

:::note
Experiments are in progress. This page will be populated with charts, data, and analysis as results become available.
:::

## What We Measured

The benchmark evaluates each routing strategy on:

- **p95 and p99 latency** (time-to-first-token and end-to-end)
- **Throughput** (requests completed per second)
- **Per-replica queue depth** over time
- **KV cache utilization** per replica

Results are reported separately for steady-state, bursty, and mixed-size workload scenarios.

---

## Steady-State Results

_Charts from Grafana/Prometheus to be inserted here._

_Key questions: Which strategy maintains the lowest tail latency under stable load? Does work distribute evenly across replicas?_

---

## Bursty Traffic Results

_Charts from Grafana/Prometheus to be inserted here._

_Key questions: How quickly does each strategy recover after a spike? Does any strategy amplify queue buildup during bursts?_

---

## Mixed Request Size Results

_Charts from Grafana/Prometheus to be inserted here._

_Key questions: Does cost-aware routing (least queue, least KV cache) provide measurable benefit when request sizes vary significantly?_

---

## Summary

_Strategy comparison table to be filled in._

| Strategy | Steady-State p99 | Bursty p99 | Mixed p99 | Notes |
|---|---|---|---|---|
| Round Robin | TBA | TBA | TBA | |
| Consistent Hashing | TBA | TBA | TBA | |
| Least Queue | TBA | TBA | TBA | |
| Least KV Cache | TBA | TBA | TBA | |

---

## Unexpected Findings

_This section will document anything that did not behave as expected — either a strategy underperformed its theoretical promise, or a simple strategy outperformed a more sophisticated one in a specific scenario._

---

## Implications for Real-World Serving

_Analysis of what the results suggest about choosing a routing strategy in production will be written here once data is available._

---

## Open Questions

_Questions raised by the results that point toward future work will be listed here._
