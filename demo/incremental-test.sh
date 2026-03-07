#!/bin/bash
# Fires 10–100 random requests every 10 seconds
# Usage: ./load_test.sh [PORT]

PORT=${1:-7778}

while true; do
    # Pick a random number of requests between 10 and 100
    NUM=$(( RANDOM % 91 + 10 ))
    echo "[$(date '+%H:%M:%S')] Sending $NUM requests..."

    for i in $(seq 1 $NUM); do
        curl -s -X POST "http://localhost:$PORT/v1/completions" \
          -H "Content-Type: application/json" \
          -d '{"model": "mistralai/Mistral-7B-v0.1", "prompt": "The following is a detailed history of computer science from the 1940s through the present day, covering key innovations in hardware, software, networking, and artificial intelligence.", "max_tokens": 1000}' > /dev/null &
    done

    echo "[$(date '+%H:%M:%S')] Burst done. Waiting 10s..."
    sleep 10
done