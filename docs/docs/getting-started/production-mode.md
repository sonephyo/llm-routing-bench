---
title: Production Mode
sidebar_position: 3
---

Production mode runs two real vLLM instances serving `mistralai/Mistral-7B-v0.1`, each pinned to a separate GPU. This is required for realistic LLM inference benchmarking.

## Requirements

- NVIDIA GPU (two GPUs recommended, one per backend)
- Docker with the NVIDIA Container Toolkit installed
- Python with `huggingface_hub` installed (`pip install huggingface_hub`)
- A Hugging Face account with access to [mistralai/Mistral-7B-v0.1](https://huggingface.co/mistralai/Mistral-7B-v0.1)

## Downloading the Model

The two vLLM backends expect the model weights to be pre-downloaded on the host machine. The `docker-compose.yml` mounts `/tmp/mistral-cache` on the host into the container at `/root/.cache/huggingface`, so download the model there:

```bash
mkdir -p /tmp/mistral-cache
HF_HOME=/tmp/mistral-cache huggingface-cli download mistralai/Mistral-7B-v0.1
```

After the download completes, your cache directory should look like:

```
/tmp/mistral-cache/
├── hub/
└── models--mistralai--Mistral-7B-v0.1/
```

:::note
You must be logged in to Hugging Face and have accepted the model's license before downloading.
Run `huggingface-cli login` and follow the prompts if you haven't already.
:::

## Docker Compose Configuration

The volume mount in `docker-compose.yml` is what connects your local model cache to the containers:

```yaml
volumes:
  - /tmp/mistral-cache:/root/.cache/huggingface
```

If you downloaded the model to a different path, update both `backend-1` and `backend-2` volume entries in `docker-compose.yml` to match. Each backend is pinned to a separate GPU (`device_ids: ['0']` and `device_ids: ['1']`) — make sure both GPUs are available before starting the stack.

## Starting the Stack

```bash
make server-up router=roundrobin
```

Replace `roundrobin` with any valid strategy: `consistanthashing` or `leastqueue`.

## Stopping the Stack

```bash
make server-down
```

## Ports

| Service | Port |
|---|---|
| Router | 7999 |
| Prometheus | 7779 |
| Grafana | 7780 |

Note: the backend vLLM instances are not exposed externally — all traffic routes through the router at `:7999`.

## Health Check

Sends a test completion request through the router to verify the full stack is up:

```bash
make health
```

This runs:

```bash
curl -s -X POST "http://localhost:7999/v1/completions" \
  -H "Content-Type: application/json" \
  -d '{"model": "mistralai/Mistral-7B-v0.1", "prompt": "The following is a detailed history of computer science", "max_tokens": 150}'
```

## Testing from a Remote Machine

:::warning Under Development
The production deployment at `indigo.cs.oswego.edu` is under active development and may not be available at all times.
:::

If the stack is running on `indigo.cs.oswego.edu`, you can send a request directly from any remote machine:

```bash
curl -s -X POST "http://indigo.cs.oswego.edu:7999/v1/completions" \
  -H "Content-Type: application/json" \
  -d '{"model": "mistralai/Mistral-7B-v0.1", "prompt": "The following is a detailed history of computer science", "max_tokens": 150}'
```

A successful response will return a JSON completion from the model routed through the load balancer.
