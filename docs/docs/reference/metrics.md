---
title: Metrics
sidebar_position: 3
---

The router exposes Prometheus metrics at `:7999/metrics`. These are scraped by the Prometheus container and visualized in Grafana.

## Exposed Metrics

| Metric | Type | Description |
|---|---|---|
| `lb_requests_total` | Counter | Total requests routed, labeled by backend |
| `lb_request_duration_seconds` | Histogram | Request latency in seconds, labeled by backend |

The `backend` label on each metric identifies which upstream server handled the request.

## Prometheus and Grafana Ports

| Mode | Prometheus | Grafana |
|---|---|---|
| Dev (`local`) | 7998 | 8000 |
| Production (`server`) | 7779 | 7780 |

Prometheus scrapes `:7999/metrics` on a fixed interval. Grafana is pre-configured to use Prometheus as its data source.
