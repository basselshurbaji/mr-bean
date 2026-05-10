#!/bin/bash

printf "\n"
read -rp "  Access token:  " TOKEN
printf "\n"

response=$(curl -s -X POST http://localhost:7489/app-token \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"app_name":"claude"}')

printf -v dashes '%58s' ''; sep="  ${dashes// /─}"

error=$(echo "$response" | grep -o '"error":"[^"]*"' | cut -d'"' -f4)
if [ -n "$error" ]; then
  printf "\n  \033[31m✗\033[0m  %s\n\n" "$error"
  exit 1
fi

token=$(echo "$response" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
app_name=$(echo "$response" | grep -o '"app_name":"[^"]*"' | cut -d'"' -f4)

printf "\n  \033[32m✓\033[0m  App token created  ·  %s\n\n" "$app_name"
printf "  \033[2mToken\033[0m\n"
printf "  %s\n\n" "$token"
printf "%s\n\n" "$sep"
printf "  \033[1mClaude Desktop\033[0m\n"
printf "  ~/Library/Application Support/Claude/claude_desktop_config.json\n\n"
printf "  \033[1mClaude Code\033[0m\n"
printf "  ~/.claude.json  (under mcpServers)\n\n"
printf '  {\n'
printf '    "mcpServers": {\n'
printf '      "mr-bean": {\n'
printf '        "command": "docker",\n'
printf '        "args": [\n'
printf '          "run", "--rm", "-i",\n'
printf '          "-e", "TOKEN=%s",\n' "$token"
printf '          "-e", "MR_BEAN_SERVER_URL=http://host.docker.internal:7489",\n'
printf '          "mr-bean-mcp:latest"\n'
printf '        ]\n'
printf '      }\n'
printf '    }\n'
printf '  }\n\n'
printf "%s\n\n" "$sep"
