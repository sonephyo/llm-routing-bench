---
title: Makefile Reference
sidebar_position: 4
---

All common operations are wrapped in `make` targets. The `router=<strategy>` variable sets the routing strategy when starting the stack.

## Targets

| Target | Description |
|---|---|
| `make server-up router=<strategy>` | Build and start the production stack (real vLLM, GPU required) |
| `make server-down` | Stop and remove the production stack |
| `make local-up router=<strategy>` | Build and start the dev stack (fake backends, no GPU) |
| `make local-down` | Stop and remove the dev stack |
| `make health` | Check that all services are reachable and healthy |
| `make local-test` | Send test requests directly to the fake backends (bypasses router) |
| `make logs` | Follow Docker Compose logs for all running services |

## Strategy Values

The `router=` variable accepts:

- `roundrobin`
- `consistanthashing`
- `leastqueue`

Example:

```bash
make local-up router=leastqueue
```
