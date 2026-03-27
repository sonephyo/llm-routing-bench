---
title: Routing Strategies
sidebar_position: 1
---

The routing strategy controls how the load balancer selects a backend for each incoming request. It is set via the `LB_STRATEGY` field in `.env` or passed as a `make` variable.

```bash
make local-up router=<strategy>
```

## Available Strategies

### `roundrobin`

Cycles through the list of backends sequentially. Each new request goes to the next backend in rotation, wrapping back to the first after the last.

Use `Mutex` for the decision of choosing backend servers.

Simple and stateless. Makes no assumptions about request cost or backend load. Works well when requests are uniform and backends are symmetric.

---

### `consistanthashing`

Uses an FNV-32a hash of the request URL to deterministically map requests to a backend on a consistent hash ring.

The same request URL will always route to the same backend, which can improve cache locality. Useful when requests with shared context (e.g. same prompt prefix or user session) benefit from hitting the same replica.

---

### `leastqueue`

Scrapes the `vllm:num_requests_running` metric from each backend at request time and routes to the backend with the fewest active requests.

This is the most load-aware strategy. It adapts dynamically when one backend falls behind.

:::note
`leastqueue` is a work in progress. Behavior may change.
:::

---

### `least-kvcache`

Read the kv cache usage of the vllm servers and route based on the lowest kv cache usage.

At high load, kv cache utilization reaches 100%, in which load balancing strategy becomes invalid.

:::note
`least-kvcache` is work in progress. More information will be added.
:::

## Setting the Strategy

**Via `.env`:**

```env
LB_STRATEGY=roundrobin
```

**Via make variable (overrides `.env`):**

```bash
make server-up router=leastqueue
```
