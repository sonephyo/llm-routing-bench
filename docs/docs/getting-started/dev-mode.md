---
title: Dev Mode
sidebar_position: 2
---

Dev mode runs the full stack using lightweight fake backends. No GPU is required. This is the recommended starting point for local development and testing without requiring GPU workloads. Utilizing dev mode will disable some load balancing functionalities that rely on llm performance (e.g. least kv cache)

## What Runs

Instead of real vLLM instances, dev mode starts two simple HTTP servers that return a fixed JSON response. Everything else — the router, Prometheus, and Grafana — runs identically to production.

## Starting the Stack

```bash
make local-up router=roundrobin
```

Replace `roundrobin` with any valid strategy: `consistanthashing` | `leastqueue` | `least-kvcache`.

## Example request

Running one request 
``` bash
curl localhost:7999
```

Running 1000 requests
```
for i in {1..1000}; do
  curl -s localhost:7999 &
done
```

## Grafana Access

After running your services, you can check the Grafana dashboard:

- **URL:** [http://localhost:8000](http://localhost:8000)  
- **Username:** `admin`  
- **Password:** `admin` 

Click the "Dashboards" option on the side bar and click on "Load Balancer & Go Runtime Metrics"

## Stopping the Stack

```bash
make local-down
```

## Ports

| Service | Port |
|---|---|
| Backend 1 | 7777 |
| Backend 2 | 7778 |
| Router | 7999 |
| Prometheus | 7998 |
| Grafana | 8000 |

This checks that all services are reachable and responding.
