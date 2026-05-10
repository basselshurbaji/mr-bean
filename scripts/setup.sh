#!/bin/bash

printf -v dashes '%58s' ''; sep="  ${dashes// /─}"

step() { printf "\n  \033[2m[%s]\033[0m  %s\n" "$1" "$2"; }
ok()   { printf "  \033[32m✓\033[0m  %s\n" "$1"; }
fail() { printf "\n  \033[31m✗\033[0m  %s\n\n" "$1"; exit 1; }

# ── Ask if user has an account ──────────────────────────────
printf "\n"
read -rp "  Do you have an account? [N/y]  " has_account
printf "\n"

if [[ "$has_account" =~ ^[Yy]$ ]]; then

  # ── Login ────────────────────────────────────────────────
  read -rp  "  Email:     " EMAIL
  read -rsp "  Password:  " PASSWORD
  printf "\n"

  step "1/2" "Logging in..."
  response=$(curl -s -X POST http://localhost:7489/auth/login \
    -H 'Content-Type: application/json' \
    -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}")

  error=$(echo "$response" | grep -o '"error":"[^"]*"' | cut -d'"' -f4)
  [ -n "$error" ] && fail "Login failed: $error"

  token=$(echo "$response" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
  ok "Logged in"

else

  # ── Register ─────────────────────────────────────────────
  read -rp  "  First name:  " FIRST
  read -rp  "  Last name:   " LAST
  read -rp  "  Email:       " EMAIL
  read -rsp "  Password:    " PASSWORD
  printf "\n"

  step "1/2" "Creating account..."
  response=$(curl -s -X POST http://localhost:7489/auth/register \
    -H 'Content-Type: application/json' \
    -d "{\"first_name\":\"$FIRST\",\"last_name\":\"$LAST\",\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}")

  error=$(echo "$response" | grep -o '"error":"[^"]*"' | cut -d'"' -f4)
  [ -n "$error" ] && fail "Registration failed: $error"

  token=$(echo "$response" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
  ok "Account created"

fi

# ── Create app token ─────────────────────────────────────────
step "2/2" "Creating app token..."
response=$(curl -s -X POST http://localhost:7489/app-token \
  -H "Authorization: Bearer $token" \
  -H 'Content-Type: application/json' \
  -d '{"app_name":"claude"}')

error=$(echo "$response" | grep -o '"error":"[^"]*"' | cut -d'"' -f4)
[ -n "$error" ] && fail "App token creation failed: $error"

app_token=$(echo "$response" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
app_name=$(echo "$response" | grep -o '"app_name":"[^"]*"' | cut -d'"' -f4)
ok "App token created  ·  $app_name"

# ── Output config ─────────────────────────────────────────────
printf "\n%s\n\n" "$sep"
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
printf '          "-e", "TOKEN=%s",\n' "$app_token"
printf '          "-e", "MR_BEAN_SERVER_URL=http://host.docker.internal:7489",\n'
printf '          "mr-bean-mcp:latest"\n'
printf '        ]\n'
printf '      }\n'
printf '    }\n'
printf '  }\n\n'
printf "%s\n\n" "$sep"
