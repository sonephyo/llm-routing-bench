#!/bin/bash
PORT=${1:-7778}

for i in $(seq 1 1000); do
    curl -s -X POST "http://129.3.20.10:$PORT/v1/completions" \
      -H "Content-Type: application/json" \
      -d '{"model": "mistralai/Mistral-7B-v0.1", "prompt": "The following is a detailed history of computer science from the 1940s through the present day, covering key innovations in hardware, software, networking, and artificial intelligence.", "max_tokens": 1000}' > /dev/null &
done

wait
echo "All 1000 requests sent"