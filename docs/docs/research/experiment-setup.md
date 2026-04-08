---
title: Experiment Setup & Methodology
sidebar_position: 3
---

## What make a good experiment?


### Possible Request Load Patterns
1) Uniform (baseline)
2) Bursty (sudden spike than quiet)
3) Ramp up (Gradual increaese)

### Possible Token Load Patterns
1) Low token (100)
2) Medium size token (1k)
3) High load token (10k)

> **Calibration note:** Token sizing for long prompts uses the approximation **1 token ≈ 4.45 characters**, measured empirically via vLLM `/tokenize` on `mistralai/Mistral-7B-v0.1` (6108 chars / 1373 tokens). Long prompts are sized to ~60% of the target token count so that input + output stays within 4096 tokens.


## Benchmark Program Diagram




