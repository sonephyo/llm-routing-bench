---
title: Prerequisites
sidebar_position: 1
---

Everything you need before running llm-routing-bench.

## Required Tools

- **Docker** and **Docker Compose** — all services run as containers.
- **Go** - to debug and run golang code locally

No other dependencies are needed for dev mode. Production mode additionally requires an NVIDIA GPU (see [Production Mode](./production-mode)).

## Environment File

Copy the example env file and fill in the two required fields:

```bash
cp .env.example .env
```

`.env` fields:

| Field | Description | Valid values |
|---|---|---|
| `MODE` | Which backend type to use | `server`, `local` |
| `LB_STRATEGY` | Routing strategy for the load balancer | `roundrobin`, `consistenthashing`, `leastqueue` |

Example:

```env
MODE=local
LB_STRATEGY=roundrobin
```

The `LB_STRATEGY` value can also be passed directly as a `make` variable, which overrides the `.env` value:

## Checkout these next

- [Dev Mode](./dev-mode) — run the stack locally without a GPU
- [Production Mode](./production-mode) — run with real vLLM on GPU hardware