MODEL = mistralai/Mistral-7B-v0.1
CACHE = /tmp/mistral-cache

.PHONY: up down logs health test clean

up:
	@docker rm -f mistral-gpu-0 mistral-gpu-1 2>/dev/null || true
	docker run -d --gpus device=0 \
		-p 7777:8000 \
		--ipc=host \
		--name mistral-gpu-0 \
		-v $(CACHE):/root/.cache/huggingface \
		vllm/vllm-openai:latest \
		--model $(MODEL)
	docker run -d --gpus device=1 \
		-p 7778:8000 \
		--ipc=host \
		--name mistral-gpu-1 \
		-v $(CACHE):/root/.cache/huggingface \
		vllm/vllm-openai:latest \
		--model $(MODEL)
	@echo "Starting servers... run 'make health' to check"

down:
	docker rm -f mistral-gpu-0 mistral-gpu-1

logs-0:
	docker logs -f mistral-gpu-0

logs-1:
	docker logs -f mistral-gpu-1

health:
	@curl -s http://localhost:7777/health && echo " Backend 1: ready" || echo " Backend 1: not ready"
	@curl -s http://localhost:7778/health && echo " Backend 2: ready" || echo " Backend 2: not ready"

test:
	@echo "Testing Backend 1..."
	@curl -s http://localhost:7777/v1/completions \
		-H "Content-Type: application/json" \
		-d '{"model": "$(MODEL)", "prompt": "Hello", "max_tokens": 10}'
	@echo "\n\nTesting Backend 2..."
	@curl -s http://localhost:7778/v1/completions \
		-H "Content-Type: application/json" \
		-d '{"model": "$(MODEL)", "prompt": "Hello", "max_tokens": 10}'