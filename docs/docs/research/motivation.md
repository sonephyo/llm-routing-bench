---
title: Motivation & Problem Statement
sidebar_position: 1
---

LLM inference is not just a machine learning problem — it is a distributed systems problem. This page explains why routing decisions at inference time matter, and what this research is trying to find out.

## The Problem

When you send a request to a large language model, something has to decide which server handles it. At small scale, it barely matters. At production scale — millions of users, dozens of GPU nodes — that decision directly affects how fast users get a response.

This research asks a pointed question: **do traditional load balancing strategies, the kind used in web servers and databases for decades, actually work well for LLM inference?**

The answer is not obvious. LLM requests behave differently from typical web traffic. A single request can take anywhere from a fraction of a second to tens of seconds depending on the output length, the model size, and the current state of the GPU. The serving infrastructure beneath it (KV cache, batching queues, memory pressure) creates hidden state that most routing algorithms ignore entirely.

## Why Tail Latency Is the Right Metric

Average latency is a comfortable number to report. Tail latency is an honest one.

For a production application serving millions of users, the **p95 and p99 latency** — the experience of the slowest 5% or 1% of requests — determines whether the product feels reliable. A routing strategy that looks good on average can still send a meaningful fraction of users into long queues behind heavy requests.

The goal of this research is to stress-test routing strategies under conditions that expose tail latency differences: variable request arrival rates, mixed workload sizes, and sustained load.

## The Gap This Research Fills

The LLM serving ecosystem has matured quickly. Frameworks like vLLM, TensorRT-LLM, and SGLang handle the inference side with significant sophistication. But **routing above those systems** — deciding which replica handles which request — has received far less attention from a distributed systems perspective.

Most deployed systems default to round-robin or a simple load metric. There is limited empirical data comparing routing strategies under controlled, production-like conditions using real LLM inference workloads.

This project provides that data.

## The Goal

To build a reproducible benchmarking testbed that measures tail latency across a range of load balancing strategies, using realistic LLM inference workloads, and to surface concrete tradeoffs that practitioners can act on when designing their serving infrastructure.

The work is applied: the output is not a new algorithm, but a clearer picture of how existing strategies behave when the workload is LLM inference rather than static file serving or database queries.
