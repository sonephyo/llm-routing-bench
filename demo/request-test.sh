#!/bin/bash

# Usage: ./benchmark.sh <num_requests> <port> <concurrent>
# Example: ./benchmark.sh 100 7777 10

NUM=${1:-10}
PORT=${2:-7778}
CONCURRENT=${3:-1}

send_request() {
    curl -s -X POST "http://localhost:$PORT/v1/completions" \
      -H "Content-Type: application/json" \
      -d '{"model": "mistralai/Mistral-7B-v0.1", "prompt": "The following is a detailed history of computer science", "max_tokens": 150}' > /dev/null
    echo "Request done"
}

export -f send_request
export PORT

seq 1 $NUM | xargs -P $CONCURRENT -I {} bash -c 'send_request'