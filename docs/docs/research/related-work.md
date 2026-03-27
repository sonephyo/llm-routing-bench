---
title: Related Work
sidebar_position: 5
---

An overview of prior work in LLM inference serving, load balancing, and distributed systems stability that this research builds on.

## LLM Inference Serving Frameworks

The primary systems this benchmark runs on top of:

**vLLM** [[paper](https://arxiv.org/abs/2309.06180)] — introduces PagedAttention, a memory management technique that virtualizes the KV cache to reduce fragmentation and improve GPU utilization. vLLM is the most widely deployed open-source LLM serving framework and is the primary inference engine used in this benchmark.

**SGLang** [[paper](https://arxiv.org/abs/2312.07104)] — a structured generation language and runtime that uses RadixAttention to share KV cache across requests with common prefixes. Relevant to this research because its caching behavior interacts directly with routing decisions around KV cache locality.

**TensorRT-LLM** — NVIDIA's inference optimization library. Focuses on kernel-level optimizations and continuous batching. Represents the high-performance closed ecosystem.

**Orca** [[paper](https://www.usenix.org/conference/osdi22/presentation/yu)] — an early and influential work on iteration-level scheduling for LLM serving. Introduced the concept of continuous batching, where the serving system does not wait for all requests in a batch to finish before starting new ones. This significantly improves GPU utilization and is now standard in most serving systems.

---

## Routing and Scheduling in LLM Serving

**Splitwise** [[paper](https://arxiv.org/abs/2311.18677)] — proposes separating the prefill and decode phases of LLM inference across different machines, since they have very different compute profiles. This is directly relevant to routing: if prefill and decode are disaggregated, the routing problem becomes a two-stage scheduling problem.

**Sarathi-Serve** [[paper](https://arxiv.org/abs/2403.02310)] — addresses the interference between prefill and decode in continuous batching systems. Shows that naive mixing of prefill and decode requests degrades tail latency, which motivates more careful scheduling at the routing level.

**DistServe** [[paper](https://arxiv.org/abs/2401.09670)] — another disaggregated serving approach that optimizes placement of prefill and decode stages to minimize latency. Highlights that routing decisions cannot be made in isolation from the serving architecture.

---

## Load Balancing in Distributed Systems

**The Power of Two Choices** [[paper](https://www.eecs.harvard.edu/~michaelm/postscripts/handbook2001.pdf)] — the foundational result showing that choosing the least loaded of two randomly selected servers, rather than one random server, reduces the maximum load from O(log n / log log n) to O(log log n). Used in the consistent hashing variant in this benchmark.

**Join-Idle-Queue** — a load balancing scheme that routes new requests to idle workers preferentially, avoiding unnecessary queuing. Has been shown to outperform round-robin under high load variance, which is characteristic of LLM workloads.

---

## Metastable Failures in Distributed Systems

**Metastability in Distributed Systems** [[paper](https://www.usenix.org/conference/osdi22/presentation/huang-lexiang)] — describes a class of failures where a system under high load reaches a stable but degraded state that persists even after the triggering load spike subsides. Relevant to LLM serving because overloaded routing decisions can create feedback loops — long queues slow throughput, which causes more queue buildup — that are difficult to recover from without shedding load.

---

## How This Work Relates

Most of the systems above focus on the **inference engine layer** — how to schedule within a single server or a disaggregated pair of servers. The routing layer above those systems, deciding which server or replica handles each incoming request, has received comparatively less empirical attention.

This research takes the inference engines as given and focuses on the routing decision, measuring how classic load balancing strategies perform when the workload is LLM inference rather than stateless HTTP requests. The findings are intended to complement the systems work above with empirical data that practitioners can use directly.

---

## References

1. Kwon et al., "Efficient Memory Management for Large Language Model Serving with PagedAttention," SOSP 2023. https://arxiv.org/abs/2309.06180
2. Zheng et al., "SGLang: Efficient Execution of Structured Language Model Programs," 2023. https://arxiv.org/abs/2312.07104
3. Yu et al., "Orca: A Distributed Serving System for Transformer-Based Generative Models," OSDI 2022. https://www.usenix.org/conference/osdi22/presentation/yu
4. Patel et al., "Splitwise: Efficient Generative LLM Inference Using Phase Splitting," 2023. https://arxiv.org/abs/2311.18677
5. Agrawal et al., "Sarathi-Serve: Efficient LLM Inference by Piggybacking Decodes with Chunked Prefills," 2024. https://arxiv.org/abs/2403.02310
6. Zhong et al., "DistServe: Disaggregating Prefill and Decoding for Goodput-Optimized Large Language Model Serving," 2024. https://arxiv.org/abs/2401.09670
7. Mitzenmacher, "The Power of Two Choices in Randomized Load Balancing," IEEE Transactions on Parallel and Distributed Systems, 2001.
8. Huang et al., "Metastable Failures in Distributed Systems," OSDI 2022. https://www.usenix.org/conference/osdi22/presentation/huang-lexiang
