# ============================================================
# llm-routing-bench
# ============================================================

MODEL = mistralai/Mistral-7B-v0.1

.PHONY: up down health test logs logs-backend-1 logs-backend-2 clean test-servers-up

health:
	@curl -s http://localhost:7777/health && echo " Backend 1: ready" || echo " Backend 1: not ready"
	@curl -s http://localhost:7778/health && echo " Backend 2: ready" || echo " Backend 2: not ready"
	@curl -s http://localhost:7779/-/healthy && echo " Prometheus: ready" || echo " Prometheus: not ready"
	@curl -s http://localhost:7780/api/health && echo " Grafana: ready" || echo " Grafana: not ready"


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

server-down:
	docker compose down

local-up:
	docker compose -f docker-compose.dev.yml build && docker compose -f docker-compose.dev.yml --env-file .env up -d

local-down:
	docker compose -f docker-compose.dev.yml down
