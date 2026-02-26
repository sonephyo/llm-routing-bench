# ============================================================
# llm-routing-bench
# ============================================================

MODEL = mistralai/Mistral-7B-v0.1

.PHONY: up down health test logs logs-backend-1 logs-backend-2 clean test-servers-up

test:
	@echo "Testing Backend 1..."
	@curl -s http://localhost:7777/v1/completions \
		-H "Content-Type: application/json" \
		-d '{"model": "$(MODEL)", "prompt": "Hello", "max_tokens": 10}'
	@echo "\n\nTesting Backend 2..."
	@curl -s http://localhost:7778/v1/completions \
		-H "Content-Type: application/json" \
		-d '{"model": "$(MODEL)", "prompt": "Hello", "max_tokens": 10}'
	@echo ""
	
logs:
	docker compose logs -f

server-up:
	DOCKER_BUILDKIT=0 docker compose -f docker-compose.yml --env-file .env up -d

local-up:
	docker compose -f docker-compose.dev.yml build && docker compose -f docker-compose.dev.yml --env-file .env up -d

local-down:
	docker compose -f docker-compose.dev.yml down
