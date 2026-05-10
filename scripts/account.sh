#!/bin/bash

printf "\n"
read -rp "  First name:  " FIRST
read -rp "  Last name:   " LAST
read -rp "  Email:       " EMAIL
read -rsp "  Password:    " PASSWORD
printf "\n"

response=$(curl -s -X POST http://localhost:7489/auth/register \
  -H 'Content-Type: application/json' \
  -d "{\"first_name\":\"$FIRST\",\"last_name\":\"$LAST\",\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}")

printf -v dashes '%58s' ''; sep="  ${dashes// /─}"

error=$(echo "$response" | grep -o '"error":"[^"]*"' | cut -d'"' -f4)
if [ -n "$error" ]; then
  printf "\n  \033[31m✗\033[0m  %s\n\n" "$error"
  exit 1
fi

token=$(echo "$response" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

printf "\n  \033[32m✓\033[0m  Account created\n\n"
printf "  \033[2mAccess token\033[0m\n"
printf "  %s\n\n" "$token"
printf "%s\n" "$sep"
printf "  To connect Claude, create an app token:  make app-token\n"
printf "%s\n\n" "$sep"
