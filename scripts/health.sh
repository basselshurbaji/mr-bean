#!/bin/bash

response=$(curl -sf http://localhost:7489/health 2>/dev/null)
if [ $? -ne 0 ]; then
  printf "\n  Service is not responding. Run 'make up' first.\n\n"
  exit 1
fi

status=$(echo "$response" | grep -o '"status":"[^"]*"' | cut -d'"' -f4)
if [ "$status" = "healthy" ]; then
  printf "\n  \033[32m✓\033[0m  Service is healthy\n\n"
else
  printf "\n  \033[31m✗\033[0m  Unexpected response: %s\n\n" "$response"
  exit 1
fi
