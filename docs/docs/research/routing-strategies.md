---
title: Routing Strategies
sidebar_position: 2
---

Four load balancing strategies are benchmarked in this project. Each represents a different philosophy about how to distribute work across inference replicas.

## The Core Question

Every routing strategy is an answer to the same question: given a new inference request and a set of available replicas, which one should handle it?

The answer depends on what information the router has, how much overhead it's willing to pay to collect that information, and what assumption it makes about how requests differ from each other.

---

## Round Robin

**The idea:** Distribute requests evenly in a fixed rotation. Replica 1, then 2, then 3, then back to 1.

Round robin is the simplest strategy and the most common default. It assumes all requests are roughly equal in cost, and that all replicas are roughly equal in capacity. Under those conditions, it does a reasonable job of spreading load.

The problem is that neither assumption holds for LLM inference. A request that generates a 10-token response is orders of magnitude cheaper than one generating 2,000 tokens. Round robin is blind to this. Over time, some replicas accumulate backlogs while others sit idle — and the router keeps sending new requests to the backlogged ones in rotation anyway.

---

## Consistent Hashing (IP-Based, with Power of Two)

**The idea:** Hash the client's IP address to select a replica. Each client is consistently routed to the same server.

Consistent hashing is widely used in caching systems because it maximizes cache reuse — the same client hits the same server, so previously computed results or loaded context may already be warm. In LLM serving, this maps to **KV cache locality**: if a user's requests share a long system prompt or conversation history, routing them to the same replica means that context may already be cached.

This implementation uses the **power of two choices** variant: instead of mapping directly to one replica, pick two candidates from the hash and route to the one with lower current load. This small change significantly reduces the chance of routing into an overloaded replica.

The tradeoff is that consistent hashing can create **hot spots**. If a small number of IP addresses generate a disproportionate share of traffic (e.g., a few enterprise customers or API aggregators), those hashes will always land on the same replicas regardless of load elsewhere.

---

## Least Queue

**The idea:** Route each request to the replica with the fewest requests currently waiting or in-flight.

Least queue is a **work-aware** strategy. It does not assume requests are equal — it tries to avoid sending work to replicas that are already backed up.

The strength of this approach is that it reacts dynamically to load imbalance. If one replica falls behind (due to a batch of heavy requests, a GC pause, or hardware variance), the router naturally stops sending it new work until the queue drains.

The limitation is visibility: the router only sees the queue depth, not the actual compute cost of the requests already in that queue. A replica with three very long-running requests may look better than one with ten short ones.

---

## Least KV Cache

**The idea:** Route to the replica with the most available KV cache memory.

This strategy is unique to LLM inference. The KV (key-value) cache is GPU memory used to store the attention state of in-progress sequences. When KV cache fills up, new requests either wait or get degraded service.

By routing toward replicas with the most available KV cache, this strategy tries to avoid the memory pressure that causes latency spikes in LLM serving. It is the most inference-aware of the four strategies — it uses a signal that is specific to how transformer models work.

The tradeoff is that KV cache availability is a proxy metric. A replica with ample cache might still be CPU or network bound. And collecting this metric requires tighter integration with the inference engine.

---

## Why These Four

These strategies span a range from **stateless and simple** (round robin) to **stateful and inference-aware** (least KV cache). Together they let us ask: does being smarter about routing actually pay off, and if so, under what conditions?

The tradeoffs for each strategy under different workload conditions are analyzed in the [Results](./results) page once experiments are complete.
