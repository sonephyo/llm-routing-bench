# ============================================================
# llm-routing-bench
# ============================================================

MODEL = mistralai/Mistral-7B-v0.1

.PHONY: up down health test logs logs-backend-1 logs-backend-2 clean test-servers-up

health:
	@echo "Sending test request to router..."
	@curl -s -X POST "http://localhost:7999/v1/completions" \
		-H "Content-Type: application/json" \
		-d '{"model": "$(MODEL)", "prompt": "The following is a detailed history of computer science", "max_tokens": 150}'

local-test:
	@echo "Testing Backend 1..."
	@curl -s http://localhost:7777/v1/completions \
		-H "Content-Type: application/json" \
		-d '{"model": "$(MODEL)", "prompt": "Hello", "max_tokens": 10}'
	@echo "\n\nTesting Backend 2..."
	@curl -s http://localhost:7778/v1/completions \
		-H "Content-Type: application/json" \
		-d '{"model": "$(MODEL)", "prompt": "Hello", "max_tokens": 10}'
	@echo ""
	@echo "Test Router..."
		@curl http://localhost:7999

logs:
	docker compose logs -f

server-up:
	DOCKER_BUILDKIT=0 docker compose -f docker-compose.yml build && DOCKER_BUILDKIT=0 $(if $(router),LB_STRATEGY=$(router),) docker compose -f docker-compose.yml --env-file .env up -d

server-down:
	docker compose down

local-up:
	docker compose -f docker-compose.dev.yml build && $(if $(router),LB_STRATEGY=$(router),) docker compose -f docker-compose.dev.yml --env-file .env up -d

local-down:
	docker compose -f docker-compose.dev.yml down
