#!/usr/bin/env bash
set -e

EXPERIMENT_IDS=(3 4 5)
GID=300
UID_VAL=42516
STRATEGIES=("roundrobin" "random" "leastqueue" "leastkvcache")
ENV_FILE="../.env"

for EXPERIMENT_ID in "${EXPERIMENT_IDS[@]}"; do
echo "====> Experiment ID: $EXPERIMENT_ID"
for strategy in "${STRATEGIES[@]}"; do
    echo "==> Running strategy: $strategy"
    cat > "$ENV_FILE" <<EOF
# Mode local or server
MODE=server
# LB_STRATEGY: random, roundrobin, leastqueue, leastkvcache
LB_STRATEGY=$strategy
EXPERIMENT_ID=$EXPERIMENT_ID
GID=$GID
UID=$UID_VAL
EOF

    docker compose down router
    docker compose up router -d

    # Wait for router to be healthy before starting bench
    echo "    Waiting for router to be ready..."
    sleep 3

    docker compose up bench -d
    echo "    Waiting for bench to finish..."
    docker compose wait bench || true

    echo "==> Done: $strategy"
done
echo "====> Experiment $EXPERIMENT_ID complete"
done
