# ============================================================
# llm-routing-bench
# ============================================================

MODEL = mistralai/Mistral-7B-v0.1

.PHONY: up down health test logs logs-backend-1 logs-backend-2 clean test-servers-up

up:
	docker compose up -d
	@echo "Waiting for backends..."
	@until curl -s http://localhost:7777/health > /dev/null 2>&1; do sleep 5; echo "  Backend 1 loading..."; done
	@echo "  Backend 1: ready"
	@until curl -s http://localhost:7778/health > /dev/null 2>&1; do sleep 5; echo "  Backend 2 loading..."; done
	@echo "  Backend 2: ready"
	@echo ""
	@echo "All services running:"
	@echo "  Backend 1:  http://localhost:7777"
	@echo "  Backend 2:  http://localhost:7778"
	@echo "  Prometheus: http://localhost:7779"
	@echo "  Grafana:    http://localhost:7780"

down:
	docker compose down

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

logs-backend-1:
	docker compose logs -f vllm-backend-1

logs-backend-2:
	docker compose logs -f vllm-backend-2

server-router-up:
	DOCKER_BUILDKIT=0 docker compose -f docker-compose.yml --env-file .env up -d

local-router-up:
	docker compose -f docker-compose.dev.yml build && docker compose -f docker-compose.dev.yml --env-file .env up -d

local-router-down:
	docker compose -f docker-compose.dev.yml down
